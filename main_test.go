package main
import ( 
	"testing"
	"fmt"
	"encoding/json"
)

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
  	if a != b {
		t.Fatal(message,fmt.Sprintf("'%v'!='%v'", a, b))	
	}	
}

func assertNotEqual(t *testing.T, a interface{}, b interface{}, message string) {
  	if a == b {
		t.Fatal(message,fmt.Sprintf("'%v'!='%v'", a, b))	
	}	
}

func TestSum(t *testing.T) {	//default test - always passes.
	total := Sum(5, 5)
	if total != 10 {
		t.Errorf("Sum was incorrect, got: %d, want: %d.", total, 10)
	}
}

func TestFullDataStr (t *testing.T) { 
	dataString := "U=229.2E+0 I=7.52E+0 Watt=1726E+0"
	nHamData := TransformData(dataString)

	 assertEqual(t,nHamData.Amps,"7.52","Unable to get Amps Data")
	 assertEqual(t,nHamData.Watts,"1726","Unable to get Watts Data")
	 assertEqual(t,nHamData.Volts,"229.2","Unable to get Volts Data")
}


func TestJSON (t *testing.T) { 	
	nHamData := &HamegData{Volts:"229.2",Amps:"7.52",Watts:"1726"}	
	var testString = "{\"volts\":\"229.2\",\"amps\":\"7.52\",\"watts\":\"1726\"}"
	
	e, err := json.Marshal(nHamData)	//JSON results
    if err != nil {
       	t.Error("Failed to JSON variable")
	}	
	
	assertEqual(t,testString,string(e),"Couldnt transform to JSON")
}

func TestBufferString(t *testing.T) { 
	var TestStr = "U=231.6E+0 I=7.59E+0 Watt=1756E+0"

	var bufferString="19,17,85,61,50,51,49,46,54,69,43,48,32,73,61,55,46,53,57,69,43,48,32,87,97,116,116,61,49,55,53,54,69,43,48,13"
	packetBytes := bufferToStringArr(bufferString)	
	calcStr := convertHexToAscii(packetBytes)		
	assertEqual(t,TestStr,calcStr,"Couldnt convert buffer to string")		
}

func TestPreSplit(t *testing.T) {
	
	splitVolts :="U=240.4E+0"		//data for splitting
	splitAmps :="I=0.01E+0"
	splitWatt :="Watt=0E+0"

	Vval := SplitMetric(splitVolts)
	Aval := SplitMetric(splitAmps)
	Wval := SplitMetric(splitWatt)

	assertEqual(t, Vval, "240.4", "Failed to split Volts")
	assertEqual(t, Aval, "0.01", "Failed to split Amps")
	assertEqual(t, Wval, "0", "Failed to split Watts")
}



// func TestNoData (t *testing.T) { 
// 	var allData="NODATA"
// 	hamegDataArr := strings.Split(allData, " ")	

	
// }
// func compareStr(expected string, got string) { 
// 	if expected!=got { 
// 		t.Error("Failed to split data, got:",got, " instead of expected:",expected)
// 	}
// }
