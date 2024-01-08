const userEl = document.getElementById("user");
const contextEl = document.getElementById("user-context");

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


