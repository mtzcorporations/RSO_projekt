// Initialize and add the map



function initMap() {
    var polyline = new google.maps.Polyline({
        // set desired options for color width
        strokeColor:"#e30f0f",  // blue (RRGGBB, R=red, G=green, B=blue)
        strokeOpacity: 0.4      // opacity of line
    });
    // The location of Uluru
    const cords = [];
    var path = [];
    const Ptuj = { lat: 1.58696889, lng: 4.6419967799 };
    const Maribor = { lat: 1.68696889, lng: 4.6419967799 };
    cords.push(Ptuj)
    cords.push(Maribor)

    // The map, centered at Uluru
    const map = new google.maps.Map(document.getElementById("map"), {
        zoom: 4,
        center: cords[0],
    });
    // Create markers
    const marker = new google.maps.Marker({
        position: cords[1],
        map: map,
    });
    const marker2 = new google.maps.Marker({
        position: cords[0],
        map: map,
    });
    path.push(marker.getPosition());
    path.push(marker2.getPosition());
    polyline.setMap(map);
    polyline.setPath(path);
}

window.initMap = initMap;