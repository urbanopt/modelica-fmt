within Somewhere;
model NdArrays
  "Example demonstrating multidimensional array wrapping (-wrap-arrays)"
  // 1D arrays stay on a single line
  parameter Integer array1D[5]={1,2,3,4,5};
  // 2D arrays break once (innermost dimension stays inline)
  parameter Integer array2D[2,2]={
    {1,2},
    {3,4}};
  // 3D arrays break once (innermost two dimensions stay inline)
  parameter Integer array3D[2,2,2]={
    {{1,2},{3,4}},
    {{5,6},{7,8}}};
  // 4D arrays break the two outer dimensions (innermost two stay inline)
  parameter Integer array4D[2,2,2,2]={
    {
      {{1,2},{3,4}},
      {{5,6},{7,8}}},
    {
      {{9,10},{11,12}},
      {{13,14},{15,16}}}};
  // iterator constructors are left untouched
  parameter Integer iter[3]={i for i in 1:3};
  // matrices keep their existing per-row formatting
  parameter Real matrixA[2,3]=[
    1.0,2.0,3.0;
    5.0,6.0,7.0];
  // arrays inside annotations remain on a single line
  parameter Real p=1.0
    annotation (Dialog(group="Fan"));
end NdArrays;
