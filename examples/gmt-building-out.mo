// Model copied from GeoJSON to Modelica Translator project (https://github.com/urbanopt/geojson-modelica-translator)
within a_project.B5a6b99ec37f4de7f94020090;
model building
  "n-zone RC building model based on URBANopt's use of TEASER export, with distribution pumps"
  extends PartialBuilding(
    redeclare package Medium=MediumW,
    have_fan=false,
    have_eleHea=false,
    have_eleCoo=false);
  package MediumW=Buildings.Media.Water
    "Source side medium";
  package MediumA=Buildings.Media.Air
    "Load side medium";
  parameter Integer nZon=6
    "Number of thermal zones";
  parameter Integer facSca=3
    "Scaling factor to be applied to on each extensive quantity";
  parameter Modelica.SIunits.TemperatureDifference delTBuiCoo=5
    "Nominal building supply and return chilled water temperature difference";
  Buildings.Controls.OBC.CDL.Continuous.Sources.Constant minTSet[nZon](
    k=fill(
      293.15,
      nZon),
    y(
      each final unit="K",
      each displayUnit="degC"))
    "Minimum temperature set point"
    annotation (Placement(transformation(extent={{-290,230},{-270,250}})));
  Buildings.Controls.OBC.CDL.Continuous.Sources.Constant maxTSet[nZon](
    k=fill(
      297.15,
      nZon),
    y(
      each final unit="K",
      each displayUnit="degC"))
    "Maximum temperature set point"
    annotation (Placement(transformation(extent={{-290,190},{-270,210}})));
  Meeting meeting
    annotation (Placement(transformation(extent={{-160,-20},{-140,0}})));
  Floor floor
    annotation (Placement(transformation(extent={{-120,-20},{-100,0}})));
  Storage storage
    annotation (Placement(transformation(extent={{-80,-20},{-60,0}})));
  Office office
    annotation (Placement(transformation(extent={{-40,-20},{-20,0}})));
  Restroom restroom
    annotation (Placement(transformation(extent={{0,-20},{20,0}})));
  ICT ict
    annotation (Placement(transformation(extent={{40,-20},{60,0}})));
  Buildings.Controls.OBC.CDL.Continuous.MultiSum mulSum(
    nin=2) if have_pum
    annotation (Placement(transformation(extent={{260,70},{280,90}})));
  Buildings.Applications.DHC.Loads.Examples.BaseClasses.FanCoil4PipeHeatPorts terUni[nZon](
    redeclare each package Medium1=MediumW,
    redeclare each package Medium2=MediumA,
    each facSca=facSca,
    QHea_flow_nominal={10000,10000,10000,10000,10000,10000},
    QCoo_flow_nominal={-10000,-10000,-10000,-10000,-10000,-50000},
    each T_aLoaHea_nominal=293.15,
    each T_aLoaCoo_nominal=297.15,
    each T_bHeaWat_nominal=35+273.15,
    each T_bChiWat_nominal=12+273.15,
    each T_aHeaWat_nominal=40+273.15,
    each T_aChiWat_nominal=7+273.15,
    each mLoaHea_flow_nominal=5,
    each mLoaCoo_flow_nominal=5)
    "Terminal unit"
    annotation (Placement(transformation(extent={{-200,-60},{-180,-40}})));
  Buildings.Applications.DHC.Loads.BaseClasses.FlowDistribution disFloHea(
    redeclare package Medium=MediumW,
    m_flow_nominal=sum(
      terUni.mHeaWat_flow_nominal .* terUni.facSca),
    dp_nominal(
      displayUnit="Pa")=100000,
    have_pum=have_pum,
    nPorts_a1=nZon,
    nPorts_b1=nZon)
    "Heating water distribution system"
    annotation (Placement(transformation(extent={{-140,-100},{-120,-80}})));
  Buildings.Applications.DHC.Loads.BaseClasses.FlowDistribution disFloCoo(
    redeclare package Medium=MediumW,
    m_flow_nominal=sum(
      terUni.mChiWat_flow_nominal .* terUni.facSca),
    typDis=Buildings.Applications.DHC.Loads.Types.DistributionType.ChilledWater,
    dp_nominal(
      displayUnit="Pa")=100000,
    have_pum=have_pum,
    nPorts_a1=nZon,
    nPorts_b1=nZon)
    "Chilled water distribution system"
    annotation (Placement(transformation(extent={{-140,-160},{-120,-140}})));
equation
  connect(disFloHea.port_b,secHeaRet[1])
    annotation (Line(points={{140,-70},{240,-70},{240,32},{300,32}},color={0,127,255}));
  connect(disFloHea.port_a,secHeaSup[1])
    annotation (Line(points={{120,-70},{-242,-70},{-242,32},{-300,32}},color={0,127,255}));
  connect(disFloCoo.port_b,secCooRet[1])
    annotation (Line(points={{140,-110},{252,-110},{252,-30},{300,-30}},color={0,127,255}));
  connect(disFloCoo.port_a,secCooSup[1])
    annotation (Line(points={{120,-110},{-280,-110},{-280,-30},{-300,-30}},color={0,127,255}));
  connect(disFloHea.ports_a1,terUni.port_bHeaWat)
    annotation (Line(points={{-120,-80.6667},{-104,-80.6667},{-104,-58.3333},{-180,-58.3333}},color={0,127,255}));
  connect(disFloHea.ports_b1,terUni.port_aHeaWat)
    annotation (Line(points={{-140,-80.6667},{-216,-80.6667},{-216,-58.3333},{-200,-58.3333}},color={0,127,255}));
  connect(disFloCoo.ports_a1,terUni.port_bChiWat)
    annotation (Line(points={{-120,-144},{-94,-144},{-94,-56},{-180,-56},{-180,-56.6667}},color={0,127,255}));
  connect(disFloCoo.ports_b1,terUni.port_aChiWat)
    annotation (Line(points={{-140,-144},{-226,-144},{-226,-56.6667},{-200,-56.6667}},color={0,127,255}));
  connect(weaBus,meeting.weaBus)
    annotation (Line(points={{1,300},{0,300},{0,20},{-66,20},{-66,-10.2},{-96,-10.2}},color={255,204,51},thickness=0.5),Text(string="%first",index=-1,extent={{6,3},{6,3}},horizontalAlignment=TextAlignment.Left));
  connect(terUni[0+1].heaPorCon,meeting.port_a)
    annotation (Line(points={{-193.333,-50},{-192,-50},{-192,0},{-90,0}},color={191,0,0}));
  connect(terUni[0+1].heaPorRad,meeting.port_a)
    annotation (Line(points={{-186.667,-50},{-90,-50},{-90,0}},color={191,0,0}));
  connect(weaBus,floor.weaBus)
    annotation (Line(points={{1,300},{0,300},{0,20},{-66,20},{-66,-10.2},{-96,-10.2}},color={255,204,51},thickness=0.5),Text(string="%first",index=-1,extent={{6,3},{6,3}},horizontalAlignment=TextAlignment.Left));
  connect(terUni[1+1].heaPorCon,floor.port_a)
    annotation (Line(points={{-193.333,-50},{-192,-50},{-192,0},{-90,0}},color={191,0,0}));
  connect(terUni[1+1].heaPorRad,floor.port_a)
    annotation (Line(points={{-186.667,-50},{-90,-50},{-90,0}},color={191,0,0}));
  connect(weaBus,storage.weaBus)
    annotation (Line(points={{1,300},{0,300},{0,20},{-66,20},{-66,-10.2},{-96,-10.2}},color={255,204,51},thickness=0.5),Text(string="%first",index=-1,extent={{6,3},{6,3}},horizontalAlignment=TextAlignment.Left));
  connect(terUni[2+1].heaPorCon,storage.port_a)
    annotation (Line(points={{-193.333,-50},{-192,-50},{-192,0},{-90,0}},color={191,0,0}));
  connect(terUni[2+1].heaPorRad,storage.port_a)
    annotation (Line(points={{-186.667,-50},{-90,-50},{-90,0}},color={191,0,0}));
  connect(weaBus,office.weaBus)
    annotation (Line(points={{1,300},{0,300},{0,20},{-66,20},{-66,-10.2},{-96,-10.2}},color={255,204,51},thickness=0.5),Text(string="%first",index=-1,extent={{6,3},{6,3}},horizontalAlignment=TextAlignment.Left));
  connect(terUni[3+1].heaPorCon,office.port_a)
    annotation (Line(points={{-193.333,-50},{-192,-50},{-192,0},{-90,0}},color={191,0,0}));
  connect(terUni[3+1].heaPorRad,office.port_a)
    annotation (Line(points={{-186.667,-50},{-90,-50},{-90,0}},color={191,0,0}));
  connect(weaBus,restroom.weaBus)
    annotation (Line(points={{1,300},{0,300},{0,20},{-66,20},{-66,-10.2},{-96,-10.2}},color={255,204,51},thickness=0.5),Text(string="%first",index=-1,extent={{6,3},{6,3}},horizontalAlignment=TextAlignment.Left));
  connect(terUni[4+1].heaPorCon,restroom.port_a)
    annotation (Line(points={{-193.333,-50},{-192,-50},{-192,0},{-90,0}},color={191,0,0}));
  connect(terUni[4+1].heaPorRad,restroom.port_a)
    annotation (Line(points={{-186.667,-50},{-90,-50},{-90,0}},color={191,0,0}));
  connect(weaBus,ict.weaBus)
    annotation (Line(points={{1,300},{0,300},{0,20},{-66,20},{-66,-10.2},{-96,-10.2}},color={255,204,51},thickness=0.5),Text(string="%first",index=-1,extent={{6,3},{6,3}},horizontalAlignment=TextAlignment.Left));
  connect(terUni[5+1].heaPorCon,ict.port_a)
    annotation (Line(points={{-193.333,-50},{-192,-50},{-192,0},{-90,0}},color={191,0,0}));
  connect(terUni[5+1].heaPorRad,ict.port_a)
    annotation (Line(points={{-186.667,-50},{-90,-50},{-90,0}},color={191,0,0}));
  connect(terUni.mReqHeaWat_flow,disFloHea.mReq_flow)
    annotation (Line(points={{-179.167,-53.3333},{-179.167,-54},{-170,-54},{-170,-94},{-141,-94}},color={0,0,127}));
  connect(terUni.mReqChiWat_flow,disFloCoo.mReq_flow)
    annotation (Line(points={{-179.167,-55},{-179.167,-56},{-172,-56},{-172,-154},{-141,-154}},color={0,0,127}));
  connect(mulSum.y,PPum)
    annotation (Line(points={{282,80},{320,80}},color={0,0,127}));
  connect(disFloHea.PPum,mulSum.u[1])
    annotation (Line(points={{-119,-98},{240,-98},{240,81},{258,81}},color={0,0,127}));
  connect(disFloCoo.PPum,mulSum.u[2])
    annotation (Line(points={{-119,-158},{240,-158},{240,79},{258,79}},color={0,0,127}));
  connect(disFloHea.QActTot_flow,QHea_flow)
    annotation (Line(points={{-119,-96},{223.5,-96},{223.5,280},{320,280}},color={0,0,127}));
  connect(disFloCoo.QActTot_flow,QCoo_flow)
    annotation (Line(points={{-119,-156},{230,-156},{230,240},{320,240}},color={0,0,127}));
  connect(maxTSet.y,terUni.TSetCoo)
    annotation (Line(points={{-268,200},{-240,200},{-240,-46.6667},{-200.833,-46.6667}},color={0,0,127}));
  connect(minTSet.y,terUni.TSetHea)
    annotation (Line(points={{-268,240},{-220,240},{-220,-45},{-200.833,-45}},color={0,0,127}));
  annotation (
    Documentation(
      info="
<html>
<p>
Building wrapper for running n-zone thermal zone models generated by TEASER.

The heating and cooling loads are computed with a four-pipe
fan coil unit model derived from
<a href=\"modelica://Buildings.Applications.DHC.Loads.BaseClasses.PartialTerminalUnit\">
Buildings.Applications.DHC.Loads.BaseClasses.PartialTerminalUnit</a>
and connected to the room model by means of heat ports.
</p>
</html>",
      revisions="<html>
<ul>
<li>
May 29, 2020, by Nicholas Long:<br/>
Template model for use in GeoJSON to Modelica Translator.
</li>
<li>
February 21, 2020, by Antoine Gautier:<br/>
First implementation.
</li>
</ul>
</html>"));
end building;
