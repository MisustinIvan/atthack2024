let app_state = {
	map: {},
	points: {},
	lines: {},
	geojson_layer: {},
};

const MARKER_RADIUS = 8;

const dataToGeoJSON = (points, edges) => {
	const geoJSON = {
		type: "FeatureCollection",
		features: [],
	};

	points.features.forEach((point) => {
		geoJSON.features.push(point);
	});

	edges.features.forEach((edge) => {
		geoJSON.features.push(edge);
	});

	return geoJSON;
};

/**
 * Custom styling for points on the map
 * @param {Object} feature - The feature object from GeoJSON
 * @param {Object} latlng - Latitude and longitude of the point
 * @returns {L.CircleMarker} - A Leaflet circle marker
 */
const pointToLayer = (point, latlng) => {
	const marker = L.circleMarker(latlng, {
		radius: MARKER_RADIUS,
		color: "blue",
		weight: 2,
		fillColor: "blue",
		fillOpacity: 0.5,
	});
	return marker;
};

/**
 * Renders GeoJSON data on the map
 * @param {L.Map} map - The Leaflet map instance
 * @param {Object} geoJSON - The GeoJSON data to render
 */
const render = (map, geoJSON) => {
	if (app_state.geojson_layer != {}) {
		app_state.map.removeLayer(app_state.geojson_layer);
	}
	app_state.geojson_layer = L.geoJSON(geoJSON, {
		pointToLayer: pointToLayer,
		onEachFeature: (feature, layer) => {
			if (feature.properties && feature.properties.name) {
				layer.bindPopup(feature.properties.name);
			}
		},
	}).addTo(map);
};

/**
 * Initializes the Leaflet map
 * @returns {L.Map} - The initialized Leaflet map
 */
const initMap = () => {
	const map = L.map("map", { doubleClickZoom: false }).setView(
		[51.505, -0.09],
		13,
	);
	L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
		maxZoom: 19,
		attribution:
			'&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>',
	}).addTo(map);
	return map;
};

let main = async () => {
	app_state.map = initMap();
	// Fetch points and parse them as JSON
	app_state.points = await fetch("/points")
		.then((response) => response.json()) // Parse response to JSON
		.catch((error) => {
			console.error("Error fetching points:", error);
			return []; // Fallback to an empty array if there's an error
		});

	app_state.lines = await fetch("/lines")
		.then((response) => response.json()) // Parse response to JSON
		.catch((error) => {
			console.error("Error fetching lines:", error);
			return []; // Fallback to an empty array if there's an error
		});

	render(app_state.map, dataToGeoJSON(app_state.points, app_state.lines));
};

document.addEventListener("DOMContentLoaded", main);
