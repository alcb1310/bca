const userEl = document.getElementById("user");
const contextEl = document.getElementById("user-context");
const opaqueEl = document.getElementById("opaque");
const drawerEl = document.getElementById("drawer");
const addButtonEl = document.getElementById("add-button");

let isContextMenuOpen = false;

userEl.addEventListener("click", (event) => {
     event.preventDefault();
     isContextMenuOpen = !isContextMenuOpen;
     contextEl.classList.toggle("show-contextmenu");
})

window.onkeyup = function(e) {
     if (isContextMenuOpen && e.keyCode === 27) {
          isContextMenuOpen = false;
          contextEl.classList.remove("show-contextmenu");
     }
}



if (opaqueEl) {
     opaqueEl.addEventListener("click", closeDrawer);
}

function openDrawer() {
     opaqueEl.style.display = "block";
     drawerEl.style.display = "flex";
}

function closeDrawer() {
     opaqueEl.style.display = "none";
     drawerEl.style.display = "none";
}
