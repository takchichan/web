const app = getApp()
Component({
    properties: {
        // defaultData 父页面（即引用组件的页面）传递的数据
        defaultData: {
            type: Object,
            value: {
                title: "我是默认标题"
            },
            observer: function(newVal, oldVal) {}
        }
    },
    data: {
        navBarHeight: app.globalData.navBarHeight,
        menuRight: app.globalData.menuRight,
        menuBotton: app.globalData.menuBotton,
        menuHeight: app.globalData.menuHeight,
    },
    attached: function() {},
    methods: {}
})