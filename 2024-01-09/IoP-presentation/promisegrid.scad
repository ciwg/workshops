$fn=100;

function rand(max) = rands(0,max,1)[0];

// gridsize = 100;
gridheight = 1;

nodeheight = 1;
minnodesize = 10;
maxnodesize = 10;
// numnodes = 10;

minappsize = 1;
maxappsize = 4;
maxnumapps = 7;

// an app is a small blue box
module app(x,y,z, u,v,w) {
    color("blue") translate([x,y,z]) cube([u,v,w]);
}

// a node is a medium-sized green box
module node(x,y,z, u,v,w) {
    color("green") translate([x,y,z]) cube([u,v,w]);
    // a node contains a random number of apps
    numapps = rand(maxnumapps);
    for (i=[1:numapps]) {
        // all sizes are between minappsize and maxappsize
        d = minappsize + rand(maxappsize-minappsize);
        e = minappsize + rand(maxappsize-minappsize);
        f = maxappsize;
        // all coordinates are between 0 and the max for that axis
        a = x+rand(u-maxappsize);
        b = y+rand(v-maxappsize);
        c = w;
        app(a,b,c,d,e,f);
    }
}

// the grid is an orange plane
module promisegrid(x,y,z, u,v,w) {
    color("orange") translate([x,y,z]) cube([u,v,w]);
}

// generate an overview of everything
module all(numnodes, gridsize) {
    // generate a large plane, color it orange, place it at 0,0,-gridheight
    promisegrid(0,0,nodeheight, gridsize,gridsize,gridheight);

    // generate numnodes nodes, each with random x,y coordinates and random sizes
    for (i=[1:numnodes]) {
        // all sizes are between minnodesize and maxnodesize
        u = minnodesize + rand(maxnodesize-minnodesize);
        v = minnodesize + rand(maxnodesize-minnodesize);
        w = nodeheight;
        // all coordinates are between 0 and gridsize-maxnodesize
        z = 0;
        // if numnodes == 1, then the node will be placed at 0,0,0
        if (numnodes == 1) {
            x = 0;
            y = 0;
            node(x,y,z,u,v,w);
        } else {
            x = rand(gridsize-maxnodesize);
            y = rand(gridsize-maxnodesize);
            node(x,y,z,u,v,w);
        }
    }

}

all(10, 100);
// all(1, 0);
