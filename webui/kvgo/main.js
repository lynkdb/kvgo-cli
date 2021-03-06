var kvgo = {
    version: "0.0.1",
    booted: false,
    base: "/kvgo/",
    api: "/kvgo/",
    _node_statuses: [
        {
            title: "LIVE",
            action: 1 << 1,
            html_tag: "success",
        },
        {
            title: "SUSPECT",
            action: 1 << 2,
            html_tag: "warning",
        },
        {
            title: "DEAD",
            action: 1 << 3,
            html_tag: "danger",
        },
    ],
    instanceEntryMenu: [
        {
            title: "Overview",
            uri: "index",
        },
        {
            title: "Metrics",
            uri: "metrics",
        },
    ],
};

kvgo.staticFilepath = function (filepath) {
    return kvgo.base + "~/" + filepath;
};

kvgo.ApiCmd = function (url, options) {
    valueui.utilx.ajax(kvgo.api + url, options);
};

kvgo.TplCmd = function (url, options) {
    valueui.utilx.ajax(kvgo.base + "~/kvgo/tpl/" + url, options);
};

kvgo.Boot = function () {
    if (kvgo.booted) {
        return;
    }
    kvgo.booted = true;
    valueui.use(["kvgo/main.css"], kvgo.load);
};

kvgo.load = function () {
    //
    var ep = valueui.newEventProxy("instances", function (instances) {
        var nav = {
            navbar_brand: {
                text: "kvgo console",
                logo_url: kvgo.staticFilepath("kvgo/img/logo-light.svg"),
            },
            navbar_nav: {
                items: [],
            },
        };

        instances.items = instances.items || [];
        kvgo.instances = instances;

        for (var i in instances.items) {
            nav.navbar_nav.items.push({
                title: instances.items[i].name,
                path: instances.items[i].name,
            });
            valueui.url.eventRegister(
                instances.items[i].name,
                kvgo.InstanceEntry,
                "valueui-layout-navbar-nav-items"
            );
        }

        valueui.layout.render("std", nav, function () {
            valueui.url.eventHandler(
                kvgo.instanceEntryDefault(),
                false,
                "valueui-layout-navbar-nav-items"
            );
        });
    });

    kvgo.ApiCmd("sys/instance-list", {
        callback: ep.done("instances"),
    });
};

kvgo.instanceEntryDefault = function (name) {
    if (!name) {
        if (kvgo.instanceActiveName) {
            name = kvgo.instanceActiveName;
        } else {
            name = valueui.sessionData.get("kvgo-console-active");
        }
    }

    if (name && kvgo.instances && kvgo.instances.items.length > 0) {
        for (var i in kvgo.instances.items) {
            if (name == kvgo.instances.items[i].name) {
                kvgo.instanceActiveName = name;
                break;
            }
        }
    }

    return kvgo.instanceActiveName;
};

kvgo.InstanceEntry = function (name) {
    console.log(name);
    name = kvgo.instanceEntryDefault(name);

    if (!name) {
        return valueui.alert.open("error", "No Instance Setup");
    }

    kvgo.instanceNavbarSetup(function () {
        valueui.url.eventHandler("index", false, "valueui-layout-module-navbar-items");
    });
};

kvgo.instanceNavbarSetup = function (cb) {
    //
    valueui.layout.moduleNavbarMenu("instance", kvgo.instanceEntryMenu);

    for (var i in kvgo.instanceEntryMenu) {
        valueui.url.eventRegister(
            kvgo.instanceEntryMenu[i].uri,
            kvgo.InstanceInvoke,
            "valueui-layout-module-navbar-items"
        );
    }

    if (cb) {
        cb();
    }
};

kvgo.InstanceInvoke = function (name) {
    // return kvgo.InstanceMetrics(); // DEBUG
    switch (name) {
        case "metrics":
            kvgo.InstanceMetrics();
            break;

        default:
            kvgo.InstanceOverview();
            break;
    }
};

kvgo.instanceNodeDataFilter = function (data) {
    data.nodes = data.nodes || [];
    data.tables = data.tables || [];

    data.node_all = {
        cpu_use: 1,
        disk_use: 0,
        disk_max: 1,
        disk_percent: 0,
        mem_use: 0,
        mem_max: 1,
        mem_percent: 0,
        status_live: 0,
        status_suspect: 0,
        status_dead: 0,
    };
    var tn = valueui.utilx.unixTimeSecond();
    for (var i in data.nodes) {
        var live_sec = tn - data.nodes[i].updated;
        if (live_sec > 3600) {
            data.nodes[i].action = 1 << 3;
        } else if (live_sec > 60) {
            data.nodes[i].action = 1 << 2;
        } else {
            data.nodes[i].action = 1 << 1;
        }

        if (!data.nodes[i].caps) {
            continue;
        }
        if (data.nodes[i].caps.disk) {
            data.node_all.disk_use += data.nodes[i].caps.disk.use;
            data.node_all.disk_max += data.nodes[i].caps.disk.max;
        }
        if (data.nodes[i].caps.mem) {
            data.node_all.mem_use += data.nodes[i].caps.mem.use;
            data.node_all.mem_max += data.nodes[i].caps.mem.max;
        }
    }
    data.node_all.disk_percent = parseInt(
        parseFloat(100 * data.node_all.disk_use) / parseFloat(data.node_all.disk_max)
    );
    if (data.node_all.disk_percent > 100) {
        data.node_all.disk_percent = 100;
    }

    data.node_all.mem_percent = parseInt(
        parseFloat(100 * data.node_all.mem_use) / parseFloat(data.node_all.mem_max)
    );

    if (data.node_all.disk_percent < 1) {
        data.node_all.disk_percent = 1;
    }

    if (data.node_all.mem_percent < 1) {
        data.node_all.mem_percent = 1;
    }

    data.node_all.status_live = data.nodes.length; // TODO

    for (var i in data.tables) {
        data.tables[i].options = data.tables[i].options || {};
        data.tables[i].key_num = data.tables[i].key_num || 0;
        data.tables[i].db_size = data.tables[i].db_size || 0;
        data.tables[i]._log_id = 0;
        data.tables[i]._options = [];
        if (data.tables[i].options) {
            for (var j in data.tables[i].options) {
                if (j == "log_id") {
                    data.tables[i]._log_id = data.tables[i].options[j];
                    continue;
                }
                data.tables[i]._options.push(j + ":" + data.tables[i].options[j]);
            }
        }
    }

    return data;
};

kvgo.InstanceOverview = function () {
    //
    var name = kvgo.instanceActiveName;

    var ep = valueui.newEventProxy("data", "tpl", function (data, tpl) {
        var msg = valueui.utilx.kindCheck(data, "SysStatus");
        if (msg) {
            return valueui.alert.open("error", msg);
        }

        valueui.sessionData.set("kvgo-console-active", name);

        data = kvgo.instanceNodeDataFilter(data.item);

        valueui.template.render({
            dstid: "valueui-layout-main",
            tplsrc: tpl,
            callback: function () {
                valueui.template.render({
                    dstid: "kvgo-instance-overview-node-list",
                    tplid: "kvgo-instance-overview-node-list-tpl",
                    data: data,
                });
                valueui.template.render({
                    dstid: "kvgo-instance-overview-node-all",
                    tplid: "kvgo-instance-overview-node-all-tpl",
                    data: data,
                });
                valueui.template.render({
                    dstid: "kvgo-instance-overview-table-list",
                    tplid: "kvgo-instance-overview-table-list-tpl",
                    data: data,
                });

                valueui.job.register({
                    id: "kvgo-instance-overview-node-list",
                    delay: 10000,
                    func: kvgo.instanceOverviewRefresh,
                });
            },
        });
    });

    ep.fail(function (err) {
        valueui.alert.open("error", "network error " + err);
    });

    kvgo.TplCmd("instance/index.htm", {
        callback: ep.done("tpl"),
    });

    kvgo.ApiCmd("sys/status?instance_name=" + name, {
        callback: ep.done("data"),
    });
};

kvgo.instanceOverviewRefresh = function (ctx) {
    var dstid = "kvgo-instance-overview";
    var elem = document.getElementById(dstid + "-node-all");
    if (!elem) {
        return ctx.callback("clean");
    }

    kvgo.ApiCmd("sys/status?instance_name=" + kvgo.instanceActiveName, {
        callback: function (err, data) {
            ctx.callback();

            if (err) {
                return;
            }
            var msg = valueui.utilx.kindCheck(data, "SysStatus");
            if (msg) {
                return;
            }

            data = kvgo.instanceNodeDataFilter(data.item);

            valueui.template.render({
                dstid: dstid + "-node-all",
                tplid: dstid + "-node-all-tpl",
                data: {
                    node_all: data.node_all,
                },
                data_hash_skip: true,
            });

            valueui.template.render({
                dstid: dstid + "-node-list",
                tplid: dstid + "-node-list-tpl",
                data: {
                    nodes: data.nodes,
                },
                data_hash_skip: true,
            });

            valueui.template.render({
                dstid: dstid + "-table-list",
                tplid: dstid + "-table-list-tpl",
                data: {
                    tables: data.tables,
                },
                data_hash_skip: true,
            });
        },
    });
};

kvgo.hchartUse = function (cb) {
    if (kvgo.hchartUsed) {
        return cb();
    }
    valueui.use(["hchart/webui/hchart.js"], function (err) {
        if (err) {
            return cb(err);
        }
        kvgo.hchartUsed = true;
        hooto_chart.basepath = kvgo.staticFilepath("hchart/webui");
        hooto_chart.opts_width = "600px";
        hooto_chart.opts_height = "400px";
        cb();
    });
};

kvgo.perfMap = {
    srk: "Read Key",
    srkr: "Read KeyRange",
    srlr: "Read LogRange",
    swk: "Write Key",
    srb: "Read Bytes",
    swb: "Write Bytes",
    //
    ark: "Read Key",
    arkr: "Read KeyRange",
    arlr: "Read LogRange",
    awk: "Write Key",
    arb: "Read Bytes",
    awb: "Write Bytes",
};

kvgo.hchartConfigTemplate = function (title) {
    return {
        type: "line",
        options: {
            title: title,
            width: "100%",
            height: "200px",
        },
        data: {
            labels: [],
            datasets: [],
        },
    };
};

kvgo.instanceMetricsDataRender = function (data, update) {
    var mrs = {
        aq: kvgo.hchartConfigTemplate("API Queries"),
        ab: kvgo.hchartConfigTemplate("API Bytes"),
        sq: kvgo.hchartConfigTemplate("Storage Queries"),
        sb: kvgo.hchartConfigTemplate("Storage Bytes"),
    };

    for (var j in data.item.keys) {
        var label = valueui.utilx.unixTimeFormat(data.item.keys[j], "i:s");
        mrs.sq.data.labels.push(label);
        mrs.sb.data.labels.push(label);
        mrs.aq.data.labels.push(label);
        mrs.ab.data.labels.push(label);
    }

    for (var i in data.item.items) {
        switch (data.item.items[i].name) {
            case "ark":
            case "arkr":
            case "arlr":
            case "awk":
                mrs.aq.data.datasets.push({
                    label: kvgo.perfMap[data.item.items[i].name],
                    data: data.item.items[i].values,
                });
                break;

            case "arb":
            case "awb":
                mrs.ab.data.datasets.push({
                    label: kvgo.perfMap[data.item.items[i].name],
                    data: data.item.items[i].values,
                });
                break;

            case "srk":
            case "srkr":
            case "srlr":
            case "swk":
                mrs.sq.data.datasets.push({
                    label: kvgo.perfMap[data.item.items[i].name],
                    data: data.item.items[i].values,
                });
                break;

            case "srb":
            case "swb":
                mrs.sb.data.datasets.push({
                    label: kvgo.perfMap[data.item.items[i].name],
                    data: data.item.items[i].values,
                });
                break;
        }
    }

    if (update !== true) {
        hooto_chart.RenderElement(mrs.aq, "kvgo-instance-metrics-api-queries");
        hooto_chart.RenderElement(mrs.ab, "kvgo-instance-metrics-api-bytes");
        hooto_chart.RenderElement(mrs.sq, "kvgo-instance-metrics-stor-queries");
        hooto_chart.RenderElement(mrs.sb, "kvgo-instance-metrics-stor-bytes");
    } else {
        hooto_chart.RenderUpdate(mrs.aq, "kvgo-instance-metrics-api-queries");
        hooto_chart.RenderUpdate(mrs.ab, "kvgo-instance-metrics-api-bytes");
        hooto_chart.RenderUpdate(mrs.sq, "kvgo-instance-metrics-stor-queries");
        hooto_chart.RenderUpdate(mrs.sb, "kvgo-instance-metrics-stor-bytes");
    }
};

kvgo.InstanceMetrics = function () {
    var name = kvgo.instanceActiveName;

    var req = "instance_name=" + name + "&time_recent=600&time_unit=10";

    var ep = valueui.newEventProxy("data", "tpl", "cerr", function (data, tpl, cerr) {
        if (cerr) {
            return valueui.alert.open("error", cerr);
        }
        var msg = valueui.utilx.kindCheck(data, "SysMetrics");
        if (msg) {
            return valueui.alert.open("error", msg);
        }

        valueui.sessionData.set("kvgo-console-active", name);

        valueui.template.render({
            dstid: "valueui-layout-main",
            tplsrc: tpl,
            callback: function () {
                kvgo.instanceMetricsDataRender(data);

                valueui.job.register({
                    id: "kvgo-instance-metrics-stor-queries",
                    delay: 3000,
                    func: kvgo.instanceMetricsRefresh,
                });
            },
        });
    });

    ep.fail(function (err) {
        valueui.alert.open("error", "network error " + err);
    });

    kvgo.hchartUse(ep.done("cerr"));

    kvgo.TplCmd("instance/metrics.htm", {
        callback: ep.done("tpl"),
    });

    kvgo.ApiCmd("sys/metrics?" + req, {
        callback: ep.done("data"),
    });
};

kvgo.instanceMetricsRefresh = function (ctx) {
    var elem = document.getElementById("kvgo-instance-metrics-stor-queries");
    if (!elem) {
        return ctx.callback("clean");
    }

    var req = "instance_name=" + kvgo.instanceActiveName + "&time_recent=60&time_unit=10";

    kvgo.ApiCmd("sys/metrics?" + req, {
        callback: function (err, data) {
            ctx.callback();

            if (err) {
                return;
            }
            var msg = valueui.utilx.kindCheck(data, "SysMetrics");
            if (msg) {
                return;
            }

            kvgo.instanceMetricsDataRender(data, true);
        },
    });
};

kvgo.NodeStatus = function (action, opts) {
    opts = opts || {};
    for (var i in kvgo._node_statuses) {
        if (kvgo._node_statuses[i].action == action) {
            var title = kvgo._node_statuses[i].title;
            if (opts.html_tag) {
                title = valueui.utilx.sprintf(
                    '<span class="badge bg-%s">%s</span>',
                    kvgo._node_statuses[i].html_tag,
                    title
                );
            }
            return title;
        }
    }
    return "";
};
