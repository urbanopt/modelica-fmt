within Somewhere;
class MyClass
  "Class to demo function definitions"
  extends Modelica.Icons.BasesPackage;

  function constructor
    "Construct to connect to a schedule in EnergyPlus"
    extends Modelica.Icons.Function;

    input Integer input1 "input 1 comment";
    output MyClass adapter;
    external "C" adapter = ExternalFunctionCall(param1, param2, param3);
  end constructor;

  function destructor "Some comment"
    extends Modelica.Icons.Function;

    input Integer input2;
    external "C" EnergyPlusInputVariableFree(input2);
  end destructor;
end MyClass;
