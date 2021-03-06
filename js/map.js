
    var platform = new H.service.Platform({
        apikey: "3AFfVpg3Do8DpOQ-lW7dRarElSXxsvt3mz4nZ1rrfU0"
    });
    var defaultLayers = platform.createDefaultLayers();
    
    //Step 2: initialize a map - this map is centered over Europe
    var map = new H.Map(document.getElementById('map'),
        defaultLayers.vector.normal.map, {
        center: { lat: 20, lng: 10 },
        zoom: 0,
        pixelRatio: window.devicePixelRatio || 1
    });
    
    //Add a resize listener to make sure that the map occupies the whole container
    window.addEventListener('resize', () => map.getViewPort().resize());
    
    //Step 3: make the map interactive
    var behavior = new H.mapevents.Behavior(new H.mapevents.MapEvents(map));
    
    // Create the default UI components
    var ui = H.ui.UI.createDefault(map, defaultLayers);
    
    //Read locations and add markers
    window.onload = function () {
        let locations = document.getElementsByClassName('cities');
        for (let i = 0; i < locations.length; i++) {
            fetch('https://geocoder.ls.hereapi.com/6.2/geocode.json?searchtext=' + locations[i].textContent + '&gen=9&apiKey=3AFfVpg3Do8DpOQ-lW7dRarElSXxsvt3mz4nZ1rrfU0')
                .then(response => response.json())
                .then(data => {
                    map.addObject(new H.map.Marker({
                        lat: data.Response.View[0].Result[0].Location.DisplayPosition.Latitude,
                        lng: data.Response.View[0].Result[0].Location.DisplayPosition.Longitude
                    }));
                })
        }
    }