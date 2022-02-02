// const agProductForm = document.getElementById("agregarProductoForm");
const linkInput = document.getElementById("linkInput");
const addButton = document.getElementById("addButton");


addButton.onclick = function(event) {
  const formData = {
    productID: ProductID,
    option: {
      Url: linkInput.value
    }
  };
  console.log(formData);
  const options = {
    method: 'POST',
    body: JSON.stringify(formData)
  };
  fetch('/agregarOpcion', options)
    .then(response => response.json())
    .then(response => {
      window.opener.pop_comunication_reciver(response);
      // Do something with response.
    });

}
