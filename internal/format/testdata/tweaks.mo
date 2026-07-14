model Tweaks "Demonstrates issue #26 output tweaks"
  parameter Real x = 1;
  Real y;
  Real z;
public
  Real w;
protected
  Real p;
initial equation
  y = pre(x);
equation
  z = der(y);
  p = sin(x) + atan2(x, y);
initial algorithm
  w := abs(x);
algorithm
  w := max(x, y);
  annotation (
    Documentation(
      revisions="<html>
<ul>
<li>Initial version.</li>
</ul>

</html>"));
end Tweaks;
