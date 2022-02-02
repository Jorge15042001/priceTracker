const popups_container_id = "popups_container";

const popups_container = document.getElementById(popups_container_id);

popups_container.onclick = function() {
  console.log(popups_container);
  const actives_popups = popups_container.children;

  actives_popups[actives_popups.length - 1].remove();
  if (popups_container.children.length == 0) popups_container.style.display = "none";
}
function load_html_in_popup(url) {
  popups_container.style.display = "flex";

  const new_popup = document.createElement("div");
  const new_object_html = document.createElement("object");

  new_popup.classList.add("popup_screen");

  new_object_html.type = "text/html";
  new_object_html.data = url;

  new_object_html.classList.add("html_window");

  new_popup.appendChild(new_object_html);

  popups_container.appendChild(new_popup);


}
function load_div_in_poppup(div_node) {
  popups_container.style.display = "flex";

  const new_popup = document.createElement("div");
  const new_object_html = div_node;

  new_popup.classList.add("popup_screen");

  new_object_html.type = "text/html";
  new_object_html.data = url;

  new_object_html.classList.add("html_window");

  new_popup.appendChild(new_object_html);

  popups_container.appendChild(new_popup);


}
