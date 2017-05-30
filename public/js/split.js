require(['split'], function(Split) {
    Split(['#tree-container', '#editor'], {
        direction: 'vertical',
        sizes: [20, 80]
    });
});