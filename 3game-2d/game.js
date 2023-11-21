let app = new PIXI.Application({ 
    width: 1024, 
    height: 768,
    transparent: true

});
document.body.appendChild(app.view);


let sprite = PIXI.Sprite.from('sample.png');
app.stage.addChild(sprite);


let elapsed = 0.0;
app.ticker.add((delta) => {
    elapsed += delta;
    sprite.x = 100.0 + Math.cos(elapsed/50.0) * 100.0;
});