var tokens = []

function pressed(c) {

    tokens[tokens.length] = c
    document.getElementById('result').innerHTML = ''

    display()
}

function display() {
    str = ""
    for (let index = 0; index < tokens.length; index++) {
        str += tokens[index];

    }
    document.getElementById('monitor').innerHTML = str

}

function getTokens() {
    str = ""
    for (let index = 0; index < tokens.length; index++) {
        str += tokens[index];

    }

    return str
}

function clearDisplay() {
    tokens = []
    document.getElementById('result').innerHTML = ''
    display()
}

function del() {
    document.getElementById('result').innerHTML = ''
    tokens.pop()
    display()
}

function cal() {
    str = ""
    for (let index = 0; index < tokens.length; index++) {
        str += tokens[index];
    }

    if (tokens.length == 0) {
        return
    }

    const xhttp = new XMLHttpRequest();

    xhttp.onload = function () {
        var statusCode = this.status
        var responseBody = this.responseText
        if (statusCode == 200) {
            if (responseBody == "NaN") {
                document.getElementById('result').style = "color: red;"
                document.getElementById('result').innerHTML = "Math Error"
                tokens = []
                return
            } else if (responseBody == "Divided By Zero") {
                document.getElementById('result').style = "color: red;"
                document.getElementById('result').innerHTML = "Math Error (" + responseBody + ")"
                tokens = []
                return
            }
            document.getElementById('result').style = "color: black;"
            document.getElementById('result').innerHTML = " = " + responseBody

            addToLog(str, responseBody)

            tokens = []
            tokens.push(responseBody)
        } else if (statusCode == 400) {
            document.getElementById('result').style = "color: red;"
            document.getElementById('result').innerHTML = "Syntax Error"
        } else {
            document.getElementById('result').innerHTML = "Unexpected Error"
        }
    }

    xhttp.open("POST", "/calculate");
    xhttp.setRequestHeader('Content-Type', 'text/plain')
    xhttp.send(str)

}

function addToLog(expression, result) {
    let log = document.createElement('div')
    log.className = 'rounded p-2 bg-light m-1'
    log.innerHTML = '<h5>' + expression + " = " + result + '</h5>'
    document.getElementById('history').appendChild(log)
    document.getElementById('log').scrollTo(0, document.getElementById('log').scrollHeight)
}

function clearLog() {
    document.getElementById('history').innerHTML = ''
}

document.addEventListener('keydown', (event) => {
    switch (event.key) {
        case "1":
            pressed("1")
            break;
        case "2":
            pressed("2")
            break;
        case "3":
            pressed("3")
            break;
        case "4":
            pressed("4")
            break;
        case "5":
            pressed("5")
            break;
        case "6":
            pressed("6")
            break;
        case "7":
            pressed("7")
            break;
        case "8":
            pressed("8")
            break;
        case "9":
            pressed("9")
            break;
        case "0":
            pressed("0")
            break;
        case "Backspace":
            del()
            break;
        case "+":
            pressed(" + ")
            break;
        case "-":
            pressed(" - ")
            break;
        case "*":
            pressed(" x ")
            break;
        case "x":
            pressed(" x ")
            break;
        case "/":
            pressed(" / ")
            break;
        case "^":
            pressed(" ^ ")
            break;
        case "(":
            pressed(" ( ")
            break;
        case ")":
            pressed(" ) ")
            break;
        case "s":
            pressed(" sqrt( ")
            break;
        case "Enter":
            cal()
            break;
        case ".":
            pressed(".")
            break;

        default:
            break;
    }
}, false);
