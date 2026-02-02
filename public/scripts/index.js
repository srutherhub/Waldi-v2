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
      swap:"outerHTML"
    });
  } catch (error) {
    console.error("Error:", error);
  }
}

function HandleHideMapSpinner() {
  document.getElementById("map-spinner").style.display="none"
}