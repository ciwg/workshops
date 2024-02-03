$fn=30;

// Orange Cube
color("orange") {
    cube([2, 2, 2], center=true);
}

// Blue Sphere
color("blue") {
    translate([0, 0, 1 + 0.5]) // Move the sphere on top of the cube. The cube is 2mm high, and we add half the sphere's diameter (0.5mm) to position it on top.
    sphere(0.5); // Sphere diameter is 1mm, so radius is 0.5mm.
}


