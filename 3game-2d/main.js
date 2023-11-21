// Electron main process (main.js)
const { app, BrowserWindow } = require('electron');

let mainWindow = null;


function createWindow() {
    mainWindow = new BrowserWindow({
        width: 1024 ,
        height: 768 ,
        useContentSize: true,
        resizable: false,
        webPreferences: {
            nodeIntegration: true, 
        },
        autoHideMenuBar: true,
    });

    mainWindow.loadFile('index.html');

    
    mainWindow.on('closed', () => {
        mainWindow = null;
    });
}

app.on('ready', createWindow);