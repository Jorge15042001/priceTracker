const optionsList = document.getElementById("optionsList");

function getStoreName(hostname) {
  if (hostname.includes("amazon")) return "amazon";
  if (hostname.includes("ebay")) return "ebay";
  if (hostname.includes("mercadolibre")) return "mercado libre";
  return "Desconocida"
}

fetch("/obtenerOpciones?ProductID=" + optionID)
  .then(response => response.json())
  .then(response => {
    document.getElementById("options-count").innerText = response.length;
    console.log(response);
    response.forEach(function(option) {

      const optID = document.createElement("p");
      optID.className += "opt-id";
      optID.innerText = "" + option.OptionID;

      const tienda = document.createElement("p");
      tienda.className += "opt-store";
      tienda.innerText = getStoreName(new URL(option.Url).hostname);


      const fechaEntrega = document.createElement("p");
      fechaEntrega.className += "opt-arrival";
      const arrivalDate = new Date(option.ArrivalDate);
      fechaEntrega.innerText = "" + arrivalDate.getDate() + "/" + (arrivalDate.getMonth() + 1) + "/" + arrivalDate.getFullYear();

      const puntuacion = document.createElement("p");
      puntuacion.className += "opt-opinion";
      puntuacion.innerText = "" + option.Opinion;

      const precio = document.createElement("p");
      precio.className += "opt-price";
      precio.innerText = "$" + option.Prices[option.Prices.length - 1].Price;

      const deleteButton = document.createElement("button");
      deleteButton.className += "opt-delete";
      deleteButton.onclick = function() {
        const data = {
          OptionID: option.OptionID
        }
        const options = {
          method: 'DELETE',
          body: JSON.stringify(data)
        };
        fetch("/eliminarOpcion", options)
          .then(response => response.json())
          .then(response => {
            console.log(response);
            window.location.reload(true);
          })

      }

      const deleteButtonIcon = document.createElement("img");
      deleteButtonIcon.src = "https://external-content.duckduckgo.com/iu/?u=https%3A%2F%2Fcdn2.iconfinder.com%2Fdata%2Ficons%2Fcleaning-19%2F30%2F30x30-10-512.png&f=1&nofb=1";

      const deleteButtonWrapper = document.createElement("div");
      deleteButtonWrapper.className += "inline-button-wrapper";

      deleteButton.appendChild(deleteButtonWrapper);
      deleteButtonWrapper.appendChild(deleteButtonIcon);


      const buyButton = document.createElement("button");
      buyButton.className += "opt-buy";
      buyButton.onclick = function() {
        window.open(option.Url);
      }

      const buyButtonIconWrapper = document.createElement("div");
      buyButtonIconWrapper.className += "inline-button-wrapper";

      const buyButtonIcon = document.createElement("img");
      buyButtonIcon.src = "https://www.freeiconspng.com/uploads/16x16-download-link-icon-11.png";

      buyButtonIconWrapper.appendChild(buyButtonIcon)
      buyButton.appendChild(buyButtonIconWrapper);

      const optWrapper = document.createElement("div");
      optWrapper.className += "optionWrapper";


      optWrapper.append(optID, tienda, fechaEntrega, puntuacion, precio, deleteButton, buyButton)

      optionsList.appendChild(optWrapper);
    })
    const button_text = document.createElement("p");
    button_text.innerText = "Agregar opcion";

    const addButton = document.createElement("button");
    addButton.className += "optionWrapper long-button"
    addButton.appendChild(button_text);

    const url_segments = window.location.pathname.split("/");
    const ProductID = url_segments[url_segments.length - 1];
    addButton.onclick = function() {
      popupForm = window.open("/agregarOpcionGUI/" + ProductID, "_blank",
        "location=no , menubar = no , status=no , titlebar=no , toolbar=no");
      popupForm.window.focus();
    }

    optionsList.appendChild(addButton);

  })

function pop_comunication_reciver(error) {
  console.log(error);
  popupForm.close();
  location.reload();
}
