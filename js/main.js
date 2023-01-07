var body = document.getElementsByTagName('body')[0];
var dark_theme_class = 'dark-theme';
var theme = getCookie('theme');

if(theme != '') {
    body.classList.add(theme);
}

document.addEventListener('DOMContentLoaded', function () {
    $('#theme-toggle').on('click', function () {

        if (body.classList.contains(dark_theme_class)) {
            body.classList.remove(dark_theme_class);
            $('#mode').text('Light Mode')
            setCookie('theme', 'light');
        }
        else {
        
            $('#mode').text('Dark Mode')
            body.classList.add(dark_theme_class);
            setCookie('theme', 'dark-theme');
        }

    });
});

// enregistrement du theme dans le cookie
function setCookie(name, value) {
    var d = new Date();
    d.setTime(d.getTime() + (365*24*60*60*1000));
    var expires = "expires=" + d.toUTCString();
    console.log(expires)
    document.cookie = name + "=" + value + ";" + expires + ";path=/";
    console.log(document.cookie)
}

// methode de recuperation du theme dans le cookie
function getCookie(cname) {
    var theme = cname + "=";
    var decodedCookie = decodeURIComponent(document.cookie);
    var ca = decodedCookie.split(';');
    for(var i = 0; i < ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(theme) == 0) {
            return c.substring(theme.length, c.length);
        }
    }
    return "";
}

const burger = document.querySelector('.burger')
const bod = document.querySelector('body');
burger.addEventListener('click', ()=>{
    burger.classList.toggle('active')
    bod.classList.toggle('open')
})


window.addEventListener("scroll",function(){
    var header = document.querySelector("header");
    header.classList.toggle("sticky",window.scrollY > 0);
})

document.getElementById("theme-toggle").addEventListener("click", function(){
    document.getElementsByTagName('body')[0].classList.toggle("dark");
})

var searchInput = document.getElementById('searchInput');
        
        searchInput.addEventListener('keyup', function() {
            var input = searchInput.value;

            if (window.location.pathname === "/search_monsters") {
                var divs = document.getElementsByClassName('monster');
            } else if (window.location.pathname === "/search_items") {
                var divs = document.getElementsByClassName('item');
            } else if (window.location.pathname === "/search_weapons") {
                var divs = document.getElementsByClassName('weapon');
            } else if (window.location.pathname === "/search_armors") {
                var divs = document.getElementsByClassName('armor');
            } else if (window.location.pathname === "/search_foods") {
                var divs = document.getElementsByClassName('food');
            } else if (window.location.pathname === "/search_skills") {
                var divs = document.getElementsByClassName('skill');
            }
            for (i = 0; i < divs.length; i++) {
                if(divs[i].getAttribute('data-title').toLocaleLowerCase().includes(input.toLocaleLowerCase())) {
                    divs[i].style.display = 'flex';
                }
                else {
                    divs[i].style.display = 'none';
                }
            }
        })