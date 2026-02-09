function deleteLoadStates() {
  document.addEventListener("htmx:afterOnLoad", (event) => {
    const indicators = document.querySelectorAll(".htmx-request");
    indicators.forEach((el) => el.classList.remove("htmx-request"));
  });

  window.onpopstate = function () {
    const indicators = document.querySelectorAll(".htmx-indicator");
    indicators.forEach((el) => {
      el.style.opacity = "0";
      el.style.visibility = "hidden";
    });
  };
}

function getBrowserLocation() {
  return new Promise((resolve, reject) => {
    navigator.geolocation.getCurrentPosition(
      (position) => {
        resolve(position.coords);
      },
      (error) => {
        reject(error);
      },
    );
  });
}

function addButtonOnClick(id, action) {
  document.getElementById(id).addEventListener("click", action);
}

function setElementValue(id, value) {
  document.getElementById(id).value = value;
}

async function myLocationHandler() {
  try {
    let coords = await getBrowserLocation();

    htmx.ajax("POST", "/api/form/browserlocation", {
      values: { lat: coords.latitude, lon: coords.longitude },
      target: "#form-address",
      swap: "outerHTML",
    });
  } catch (error) {
    console.error("Error:", error);
  }
}

function HandleHideMapSpinner() {
  document.getElementById("map-spinner").style.display = "none";
}

function addAnnotationstoMap() {
  const elements = document.querySelectorAll(".map-place");

  const annotations = Array.from(elements).map((el) => {
    const lat = parseFloat(el.dataset.lat);
    const lon = parseFloat(el.dataset.lon);
    const name = el.dataset.name;

    const coord = new mapkit.Coordinate(lat, lon);

    const anno = new mapkit.MarkerAnnotation(coord, {
      title: name,
      color: "#F2B705",
      enabled: false,
      size: { width: 32, height: 32 },
    });

    return anno;
  });

  map.addAnnotations(annotations);
}

async function openPlaceDetails(event) {
  const dialog = document.getElementById("map-modal");
  const elem = event.currentTarget;

  const id = elem.id;
  const lat = elem.dataset.lat;
  const lon = elem.dataset.lon;

  const contentDiv = document.getElementById("map-modal-content");
  if (contentDiv) {
    contentDiv.innerHTML = "<p>Loading...</p>";
  }

  try {
    htmx.ajax("POST", "/api/form/mapmodal", {
      values: { id: id, lat: lat, lon: lon, _t: Date.now() },
      target: "#map-modal-content",
      swap: "innerHTML",
    });
    console.log("Hello from modal");
    dialog.showModal();
  } catch (error) {
    console.error("Error:", error);
  }
}
