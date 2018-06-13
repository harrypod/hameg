https://podtech.io
------------------------------------------------------------------------------------------------------
For interogating a HM8143 Three-Channel Arbitrary Power Supply over a USB or Serial connection.

Calling only the VAL? command, this can be improved with additional commands as needed.
------------------------------------------------------------------------------------------------------

./main.exe --interface 6 --os 1 --log 2      //execute on COM6 with INFO.

Returns a JSON string:

{"volts":"239.1","amps":"0.01","watts":"0"}

------------------------------------------------------------------------------------------------------

Usage of:
  -interface int
        COM or TTY of Interface
  -log int
        Log level, 0=Off,1=Error,2=Info,3=Warning,4=Trace
  -os int
        Operating System, 0=Linux,1=Windows