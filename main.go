// ./main.exe --interface 6 --os 1 --log 2

package main
import (
 "flag"					// command line parameters
  "os"					// command line parameters + logger definitions
  "fmt"
  "time"
  "bytes"					// handling buffer
  "strconv"					// string conversion
  "strings"				
  "io"
  "log"
  "encoding/json"
  "github.com/goburrow/serial"   
  "io/ioutil"			// Logging
)

var (
    Trace   *log.Logger
    Info    *log.Logger
    Warning *log.Logger
    Error   *log.Logger
)

type HamegData struct { 
	Volts string `json:"volts"`
	Amps string  `json:"amps"`
	Watts string  `json:"watts"`	
}

func Init(traceHandle io.Writer,infoHandle io.Writer, warningHandle io.Writer,errorHandle io.Writer) {
	Trace = log.New(traceHandle,"TRACE: ",log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(infoHandle,"INFO: ",log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(warningHandle,"WARNING: ",log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorHandle,"ERROR: ",log.Ldate|log.Ltime|log.Lshortfile)
}

// setup constants for connectivity and commands sent to device.
const (  
	PARITY		= "N"		// needs to be set on LINUX
	BAUDRATE	= 9600
	STOPBITS	= 1
	DATABITS	= 8
	TIMEOUT		= 250 * time.Millisecond 
	

	cmdGetVersion 	= "\x56\x45\x52\x53\x49\x4F\x4E\x3F\x0D"		//VERSION?<CR>
	cmdGetVal 		= "\x56\x41\x4C\x3F\x0D" 						//VAL?<CR>			
)



func main() {	
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)	
	var osPtr = flag.Int("os", 0, "Operating System, 0=Linux,1=Windows")	
	var serialInterfacePtr = flag.Int("interface", 0, "COM or TTY of Interface")			// AKA INTERFACE
	var logPtr = flag.Int("log", 0, "Log level, 0=Off,1=Error,2=Info,3=Warning,4=Trace")	
	var osInterface = "/dev/ttyUSB"	

	flag.Parse()	
	switch *logPtr {            /// pass rcx as int
		case 0: // NO LOGS
			Init(ioutil.Discard,ioutil.Discard,ioutil.Discard,ioutil.Discard)
        case 1:                             // Show Error log only
            Init(ioutil.Discard,ioutil.Discard,ioutil.Discard,os.Stderr)
        case 2:                             // Info = INFO + WArning + Error
            Init(ioutil.Discard,os.Stdout,os.Stdout,os.Stderr)
        case 3:                             // Warning = Warning + Error only
            Init(ioutil.Discard,ioutil.Discard,os.Stdout,os.Stderr)
        case 4:
            Init(os.Stdout,os.Stdout,os.Stdout,os.Stderr)
        default:
            Error.Println("Invalid Log, Setting 0")
	}
		
	start := time.Now()	
	if *osPtr == 1 {					// serial interface specified by OS and COM, so COM4, or /dev/ttyUSB0 
		osInterface = "COM"
	}
	var serialInterface = osInterface + strconv.Itoa(*serialInterfacePtr)		

	Info.Println("Connecting to Serial Port: ",serialInterface)	
	port, err := serial.Open(&serial.Config{Parity: PARITY,Address: serialInterface, BaudRate: BAUDRATE, DataBits: DATABITS, StopBits: STOPBITS, Timeout: TIMEOUT})		
	if err != nil {
		Error.Printf("Could not open connection to port, check interface and OS: %v", err)
		time.Sleep(3 * time.Second)					/// sleep for 1millisecond		
	}				
	defer port.Close()				// close port on exit main
	
	allData := transmitHandle(port,[]byte(cmdGetVal))	// send cmd and return data 

	if(allData != "NODATA") { 		
		nHamData := TransformData(allData)
		e, err := json.Marshal(nHamData)	//JSON results
		if err != nil {
			Error.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(e))				//output json to stdout
	}else { 
		Error.Println("No data received")
		fmt.Println("No data received")
	}
	
	elapsed := time.Since(start).String()
	Info.Println("Time Elapsed:",elapsed)	
}

//TransformData breaks down the string, first pass into 3 metrics
func TransformData(allData string) *HamegData{ 
	hamegDataArr := strings.Split(allData, " ")			// put returned data into array		
	nHamData := &HamegData{Volts:SplitMetric(hamegDataArr[0]),Amps:SplitMetric(hamegDataArr[1]),Watts:SplitMetric(hamegDataArr[2])}
	return nHamData
}

//SplitMetric removes the metric element ofthe string and maintains the value
func SplitMetric(in string)string { 					// handle format of data provided by Hameg
	var valuesArr = strings.Split(in,"=") 				// split data by '='
	var activeValue = valuesArr[1]						// remove E+0 from each field
	var removeEPlus = activeValue[:len(activeValue)-3]
	return removeEPlus
}

func convertInttoASCII (iArr[] int) string { 
	//var combineString
	var buffer bytes.Buffer
	for _, i := range iArr { 
		buffer.WriteString(string(i))		
	}
 	return buffer.String()
}

func strToint(s string) int { 		// convert string to integer
	var i int
	i,_ = strconv.Atoi(s)
	return i
}

func strArrtoInt (s[] string) []int { 		// convert array of string to array of integer
	var retArr = []int{}
	for _, i := range s {		
        retArr = append(retArr, strToint(i))		// call strToInt func
	}	
	return retArr
}

func stringToASCII (s[] string) string { 
	var intArr = strArrtoInt(s) 				// convert string array to int array
	return convertInttoASCII(intArr)			// pass int array and return ascii in single string

}

func transmitHandle (sport serial.Port,trx []byte)string {				// send command and process data	
	Trace.Println("TRX:",trx)
	var buffer bytes.Buffer						// buffer for response
	
	n, err := sport.Write([]byte(trx))			// write decoded string to the open port.
	if err != nil {
		Error.Println("port.Write: ", err, " with trx ",n," bytes")
	}	
 	time.Sleep(TIMEOUT)							//timeout required for Hameg to provide response
	var rcxPasses =0
	for {
		rcxPasses++
		buf := make([]byte, 128)
		n, err := sport.Read(buf)
		if err != nil {
			if err != io.EOF {
				// on windows - this occurs for every node not responding				
				break
			}
		} else {
			buf = buf[:n]			
			if(len(buf)>0) { 					// GOT DATA
				if(rcxPasses >1) { 				// looped more than once on receiving data
					buffer.WriteString(",")		// add delimeter to next pass
				}					
				var rcxBuf = delimitBytes(buf,",");				// put bytes into a delimitered string										
				buffer.WriteString(rcxBuf)						// put data into rcx_buffer, need to get full packet before processing												
			}
		}
		if(len(buf) ==0 ) { 
			break									// no data , return
		}else { 
			if len(buf) > 0 {
				if buf[len(buf)-1] == 4 {			// 04 indicates end of packet
					break							// got data, return
				}
			}
		}		
	}
	if buffer.String() != "" {  								// got a buffer of data				
		packetBytes := bufferToStringArr(buffer.String())	
		return convertHexToAscii(packetBytes)						//return data after removing <CR>
	}	
	return "NODATA"				// bad handle - no node exists
}
func bufferToStringArr(bufferString string) []string { 		
	packetStrArr := strings.Split(bufferString, ",")			// convert and split string into string array	
	return packetStrArr
}


func convertHexToAscii(fPacket []string) string { 	// convert array of hex to readable string
	var data = fPacket[:len(fPacket)-1]				// remove <CR>
	var dataStr = stringToASCII(data)
	return dataStr
}

func delimitBytes(a []byte, delim string) string {				
    return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")				// CONVERT BYTE ARRAY to String Delimited    
}

//Sum is for testing only
func Sum(x int, y int) int {
    return x + y
}
