model Example
  "A short description"
  parameter Real x=1
    "a normal string comment";
  parameter String s="not html, just plain text";
  annotation (
    Documentation(
      info="
        <html>
          <p>
            Hello &amp; welcome to <b>Example</b>.
          </p>
          <p>
            See <a href=\"modelica://Some.Package.Example\">Example</a> for more details and a four-pipe fan coil unit model derived from the base classes.
          </p>
          <pre>
  code line one
    indented code line
</pre>
        </html>",
      revisions="
        <html>
          <ul>
            <li>
              January 1, 2020, by Someone:<br/>Initial implementation.
            </li>
          </ul>
        </html>"),
    __partial="<html><p>partial without a closing tag",
    experiment(
      StopTime=1));
end Example;
