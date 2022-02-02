const productList = document.getElementById("productList");
let popupForm;

fetch('/obtenerProductos')
  .then(response => response.json())
  .then(response => {

    document.getElementById("product-count").innerText = "" + response.length;

    response.forEach(product => {

      console.log(product.ProductID);

      let productImage = document.createElement("img");
      productImage.src = getPhoto(product);

      const productName = document.createElement("p");
      productName.innerText = product.ProductName;

      const productPrice = document.createElement("p");
      productPrice.innerText = getMinPrice(product);


      const imageWrapper = document.createElement("div");
      imageWrapper.className += "imageWrapper"
      imageWrapper.appendChild(productImage);

      const nameWrapper = document.createElement("div");
      nameWrapper.className += "nameWrapper";
      nameWrapper.appendChild(productName);

      const priceWrapper = document.createElement("div");
      priceWrapper.className += "priceWrapper";
      priceWrapper.appendChild(productPrice);

      const closeButton = document.createElement("button");
      closeButton.className += "conerButton";
      closeButton.innerText = "x";
      closeButton.onclick = function() {
        const data = {
          ProductID: product.ProductID
        }
        const options = {
          method: 'DELETE',
          body: JSON.stringify(data)
        };
        fetch('/eliminarProducto', options)
          .then(response => response.json())
          .then(response => {
            console.log(response);
            window.location.reload(true);
            // Do something with response.
          });

      }

      const productWrapper = document.createElement("div");
      productWrapper.className += "productWrapper";
      productWrapper.append(closeButton, imageWrapper, nameWrapper, priceWrapper);

      const redirect = function() {
        window.location = "/producto/" + product.ProductID;
      }
      imageWrapper.onclick = redirect
      nameWrapper.onclick = redirect
      priceWrapper.onclick = redirect


      productList.appendChild(productWrapper);
    });
    //add plus button

    const addButton = document.createElement("button");
    addButton.className += "full-button";
    addButton.innerText = "+";

    const addButtonWrapper = document.createElement("div");
    addButtonWrapper.className += "productWrapper";
    addButtonWrapper.appendChild(addButton);
    productList.appendChild(addButtonWrapper);

    addButton.onclick = function() {
      popupForm = window.open("/agregarProductoGUI", "_blank",
        "location=no , menubar = no , status=no , titlebar=no , toolbar=no");
      popupForm.window.focus();
    }


  });

function pop_comunication_reciver(error) {
  console.log(error);
  popupForm.close();
  location.reload();
}

function getMinPrice(product) {
  if (product.Options.length == 0) return "$0.00";

  const prices = product.Options.map((option) => {
    if (option.Prices.length == 0) return 0.0;
    return option.Prices[option.Prices.length - 1].Price / 1.0;
  });
  console.log(prices);
  // prices.forEach(p => console.log("type", typeof (p)))
  minPrice = Math.min(...prices);
  console.log(minPrice)
  return "$" + minPrice;
}

function getPhoto(product) {
  if (product.Options.length == 0) return "https://tacm.com/wp-content/uploads/2018/01/no-image-available.jpeg";
  return product.Options[0].ImageSrc;
}

