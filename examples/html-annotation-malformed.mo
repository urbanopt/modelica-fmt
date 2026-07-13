model MalformedExample "Fixture with malformed embedded HTML"
  annotation (
    Documentation(
      info="<html>
<p>This paragraph is never closed.
</html>"),
    experiment(StopTime=1));
end MalformedExample;
