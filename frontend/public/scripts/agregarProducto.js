// const agProductForm = document.getElementById("agregarProductoForm");
const nameInput = document.getElementById("nameInput");
const addButton = document.getElementById("addButton");


addButton.onclick = function(event) {
  const formData = {
    ProductName: nameInput.value,
    TrackingInterval: 60 * 60 * 24 //evey hour
  };
  const options = {
    method: 'POST',
    body: JSON.stringify(formData)
  };
  fetch('/agregarProducto', options)
    .then(response => response.json())
    .then(response => {
      window.opener.pop_comunication_reciver(response);
      // Do something with response.
    });

}
