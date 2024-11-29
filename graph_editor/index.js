const Status = {
	CREATING: "creating",
	NOT_CREATING: "not_creating",
	CREATING_ROOT: "creating_root",
};

const MARKER_RADIUS = 8;

let app_state = {
	next_id: 0,
	status: Status.NOT_CREATING,
	map: {},
	from_point: [],
	points: [],
	edges: [],
	graph: {},
	geojson_layer: {},
};

/**
 * @param{Array<number>}p1
 * @param{Array<number>}p2
 * @returns{boolean}
 **/
let pos_equals = (p1, p2) => {
	return p1[0] == p2[0] && p1[1] == p2[1];
};

/**
 * @param{Array<Array<number>>}points
 * @param{Array<number>}point
 * @returns{boolean}
 **/
let points_contain = (points, point) => {
	for (let other_point of points) {
		if (pos_equals(other_point, point)) {
			return true;
		}
	}
	return false;
};

/**
 * @param{Array<Array<Array<number>>>}points
 * @param{Array<Array<number>>}point
 * @returns{boolean}
 **/
let edges_contain = (edges, edge) => {
	return edges.some(
		(e) =>
			(e[0][0] === edge[0][0] &&
				e[0][1] === edge[0][1] &&
				e[1][0] === edge[1][0] &&
				e[1][1] === edge[1][1]) ||
			(e[0][0] === edge[1][0] &&
				e[0][1] === edge[1][1] &&
				e[1][0] === edge[0][0] &&
				e[1][1] === edge[0][1]),
	);
};

let undo = (_) => {
	console.log("undo");
	if (app_state.points.length <= 1) {
		console.log("empty points");
		console.log(app_state.points);
		console.log(app_state.edges);
		return;
	}

	app_state.points.pop();
	app_state.edges.pop();
	app_state.from_point = app_state.points[app_state.points.length - 1];
	app_state.graph = dataToGeoJSON(app_state.points, app_state.edges);
	app_state.map.removeLayer(app_state.geojson_layer);
	render(app_state.map, app_state.graph);
};

let to_save_format = (graph) => {
	let points = [];
	let lines = [];

	for (let feature of graph.features) {
		if (feature.geometry.type == "Point") {
			feature.properties = {
				id: app_state.next_id++,
			};
			points.push(feature);
		}
		if (feature.geometry.type == "LineString") {
			feature.properties = {
				state: 0,
				size: 1,
				cars: 0,
			};
			lines.push(feature);
		}
	}

	return [
		{ type: "FeatureCollection", features: points },
		{ type: "FeatureCollection", features: lines },
	];
};

let save = (_) => {
	console.log("saving");
	//console.log(app_state.graph);
	data = to_save_format(app_state.graph);
	console.log(data);
	fetch("/save_geojson", {
		method: "POST",
		body: JSON.stringify(data),
	});
};

let rerender = (_) => {
	console.log("rerender");
	console.log(app_state.points);
	console.log(app_state.edges);
	app_state.map.removeLayer(app_state.geojson_layer);
	app_state.graph = dataToGeoJSON(app_state.points, app_state.edges);
	render(app_state.map, app_state.graph);
};

/**
 * Converts a list of points and edges into GeoJSON format
 * @param {number[][]} points - List of points as [longitude, latitude] pairs
 * @param {number[][]} edges - List of edges as [point1, point2] pairs
 * @returns {Object} GeoJSON FeatureCollection object
 */
const dataToGeoJSON = (points, edges) => {
	const geoJSON = {
		type: "FeatureCollection",
		features: [],
	};

	points.forEach((point) => {
		geoJSON.features.push({
			type: "Feature",
			geometry: {
				type: "Point",
				coordinates: point,
			},
		});
	});

	edges.forEach((edge) => {
		geoJSON.features.push({
			type: "Feature",
			geometry: {
				type: "LineString",
				coordinates: edge,
			},
		});
	});

	return geoJSON;
};

/**
 * Converts the map event object to a point in [lng, lat] format
 * @param {Object} ev - The event object from map click
 * @returns {number[]} - The point as [longitude, latitude]
 */
const evToPoint = (ev) => {
	return [ev.latlng.lng, ev.latlng.lat];
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

	marker.on("click", (_) => {
		app_state.edges.push([
			app_state.from_point,
			point.geometry.coordinates,
		]);
		app_state.graph = dataToGeoJSON(app_state.points, app_state.edges);
		app_state.from_point = point.geometry.coordinates;
		render(app_state.map, app_state.graph);
	});

	marker.on("contextmenu", (_) => {
		console.log("Selected node");
		app_state.from_point = point.geometry.coordinates;
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

	console.log(dataToGeoJSON(app_state.points, app_state.edges));
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

/**
 * Main application logic
 */
const main = () => {
	const map = initMap();
	app_state.map = map;

	map.on("click", (ev) => {
		switch (app_state.status) {
			case Status.NOT_CREATING:
				return;

			case Status.CREATING_ROOT:
				console.log("Creating root...");
				app_state.status = Status.CREATING;
				app_state.from_point = evToPoint(ev);
				app_state.points.push(app_state.from_point);
				app_state.graph = dataToGeoJSON(
					app_state.points,
					app_state.edges,
				);
				render(map, app_state.graph);
				break;

			case Status.CREATING:
				console.log("Creating edge...");
				const newPoint = evToPoint(ev);
				if (!points_contain(app_state.points, newPoint)) {
					app_state.points.push(newPoint);
					app_state.edges.push([app_state.from_point, newPoint]);
					app_state.from_point = newPoint;
				}

				app_state.graph = dataToGeoJSON(
					app_state.points,
					app_state.edges,
				);
				render(map, app_state.graph);
				break;
		}
	});
};

document.addEventListener("DOMContentLoaded", main);

/**
 * Handle custom request headers for HTMX requests
 */
document.addEventListener("htmx:configRequest", (evt) => {
	const auth = "ligmaballs"; // Just a placeholder for example
	evt.detail.headers["Authorization"] = auth;
});

/**
 * Handle state changes based on HTMX request URL
 */
document.addEventListener("htmx:configRequest", (event) => {
	const url = event.detail.path;
	if (url === "/create_graph") {
		console.log("Posting to /create_graph route");
		app_state.status = Status.CREATING_ROOT;
	} else if (url === "/save_graph") {
		console.log("Posting to /stop_creating_graph route");
		app_state.status = Status.NOT_CREATING;
		save();
	}
});

document.addEventListener("htmx:afterSwap", (ev) => {
	if (
		ev.detail.target.id == "creator_controls" &&
		app_state.status == Status.CREATING_ROOT
	) {
		htmx.on("#undo", "click", undo);
		htmx.on("#rerender", "click", rerender);
	}
});
