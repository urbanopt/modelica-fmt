within Somewhere;
class MyClass
  "Class to demo formatting arrays and matrices"
  extends Modelica.Icons.BasesPackage;

  parameter Real var1=3.14;

  parameter Real x[3];

  parameter Real y[3]={1.0,0.0,-1.0};

  parameter Real z[5]={1.0,0.0,-1.0,2.0,0.0};

  parameter Real A[2,3]={{1.0,2.0,3.0},{5.0,6.0,7.0}};

  parameter Real B[:,3]={{1.0,2.0,3.0},{5.0,6.0,7.0},{1.0,2.0,3.0},{5.0,6.0,7.0},{1.0,2.0,3.0},{5.0,6.0,7.0},{1.0,2.0,3.0},{5.0,6.0,7.0}};

  parameter Real C[2,3]=[
    1.0,2.0,3.0;
    5.0,6.0,7.0];

  parameter Integer D[4,3]=[
    0,1,1;
    2,3,5;
    8,13,21;
    34,55,89];

  parameter Real fraPFan_nominal(
    unit="W/(kg/s)")=275/0.15
    "Fan power divided by water mass flow rate at design condition"
    annotation (Dialog(group="Fan"));

  parameter Modelica.SIunits.Power PFan_nominal=fraPFan_nominal*m_flow_nominal
    "Fan power"
    annotation (Dialog(group="Fan"));

end MyClass;
