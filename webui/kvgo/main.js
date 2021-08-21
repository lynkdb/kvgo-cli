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
        data.tables[i].states = data.tables[i].states || {};
        data.tables[i].key_num = data.tables[i].key_num || 0;
        data.tables[i].db_size = data.tables[i].db_size || 0;
        data.tables[i]._log_id = 0;
        data.tables[i]._states = [];
        for (var j in data.tables[i].states) {
            if (j == "log_id") {
                data.tables[i]._log_id = data.tables[i].states[j];
                continue;
            }
            data.tables[i]._states.push({
                key: j,
                value: data.tables[i].states[j],
            });
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

kvgo.hchartConfigTemplate = function (title) {
    return {
        type: "line",
        options: {
            title: title,
            width: "100%",
            height: "200px",
            radius: 0,
        },
        data: {
            labels: [],
            datasets: [],
        },
    };
};

kvgo.instanceMetricsDataRender = function (data, update) {

    if (!data.item || !data.item.metrics) {
        return
    }
    var mrs = {
        service_qps: kvgo.hchartConfigTemplate("Service Queries/10s"),
        service_siz: kvgo.hchartConfigTemplate("Service Bytes/10s"),
        service_lat: kvgo.hchartConfigTemplate("Service Latency (μs)"),
        storage_qps: kvgo.hchartConfigTemplate("Storage Queries/10s"),
        storage_siz: kvgo.hchartConfigTemplate("Storage Bytes/10s"),
        storage_lat: kvgo.hchartConfigTemplate("Storage Latency (μs)"),
        logsync_qps: kvgo.hchartConfigTemplate("LogSync Queries/10s"),
        logsync_siz: kvgo.hchartConfigTemplate("LogSync Bytes/10s"),
        logsync_lat: kvgo.hchartConfigTemplate("LogSync Latency (ms)"),
        system_cpu: kvgo.hchartConfigTemplate("CPU load (%)"),
        system_mem: kvgo.hchartConfigTemplate("Memory Usage (MiB)"),
        system_io: kvgo.hchartConfigTemplate("IO Status (MiB)"),
    };

    for (var j in data.item.time_buckets) {
        var label = valueui.utilx.unixTimeFormat(data.item.time_buckets[j], "i:s");
        for (var m in mrs) {
            mrs[m].data.labels.push(label);
        }
    }

    for (var i in data.item.metrics) {
        var m = data.item.metrics[i];
        if (!m.labels) {
            m.labels = [];
        }
        var mTPL = {
            label: null,
            data: [],
        }
        for (var k in m.labels) {
            if (mTPL.label) {
                mTPL.label += "/";
            } else {
                mTPL.label = "";
            }
            mTPL.label += m.labels[k].name + "/" + m.labels[k].value;
        }
        if (!mTPL.label) {
            mTPL.label = m.name;
        }

        switch (m.name) {
            case "ServiceCall":
                var mQPS = valueui.utilx.objectClone(mTPL);
                var mSiz = valueui.utilx.objectClone(mTPL);
                for (var j in m.points) {
                    if (m.points[j].count > 0) {
                        mQPS.data.push(m.points[j].count / 10);
                    } else {
                        mQPS.data.push(0);
                    }
                    if (m.points[j].sum > 0) {
                        mSiz.data.push(m.points[j].sum / 10);
                    } else {
                        mSiz.data.push(0);
                    }
                }
                mrs.service_qps.data.datasets.push(mQPS);
                mrs.service_siz.data.datasets.push(mSiz);
                break;

            case "ServiceLatency":
                var mLat = valueui.utilx.objectClone(mTPL);
                for (var j in m.points) {
                    if (m.points[j].count > 0) {
                        mLat.data.push(m.points[j].sum / m.points[j].count);
                    } else {
                        mLat.data.push(0);
                    }
                }
                mrs.service_lat.data.datasets.push(mLat);
                break;

            case "StorageCall":
                var mQPS = valueui.utilx.objectClone(mTPL);
                var mSiz = valueui.utilx.objectClone(mTPL);
                for (var j in m.points) {
                    if (m.points[j].count > 0) {
                        mQPS.data.push(m.points[j].count / 10);
                    } else {
                        mQPS.data.push(0);
                    }
                    if (m.points[j].sum > 0) {
                        mSiz.data.push(m.points[j].sum / 10);
                    } else {
                        mSiz.data.push(0);
                    }
                }
                mrs.storage_qps.data.datasets.push(mQPS);
                mrs.storage_siz.data.datasets.push(mSiz);
                break;

            case "StorageLatency":
                var mLat = valueui.utilx.objectClone(mTPL);
                for (var j in m.points) {
                    if (m.points[j].count > 0) {
                        mLat.data.push(m.points[j].sum / m.points[j].count);
                    } else {
                        mLat.data.push(0);
                    }
                }
                mrs.storage_lat.data.datasets.push(mLat);
                break;

            case "LogSyncCall":
                var mQPS = valueui.utilx.objectClone(mTPL);
                var mSiz = valueui.utilx.objectClone(mTPL);
                for (var j in m.points) {
                    if (m.points[j].count > 0) {
                        mQPS.data.push(m.points[j].count / 10);
                    } else {
                        mQPS.data.push(0);
                    }
                    if (m.points[j].sum > 0) {
                        mSiz.data.push(m.points[j].sum / 10);
                    } else {
                        mSiz.data.push(0);
                    }
                }
                mrs.logsync_qps.data.datasets.push(mQPS);
                mrs.logsync_siz.data.datasets.push(mSiz);
                break;

            case "LogSyncLatency":
                var mLat = valueui.utilx.objectClone(mTPL);
                for (var j in m.points) {
                    if (m.points[j].count > 0) {
                        mLat.data.push(m.points[j].sum / m.points[j].count);
                    } else {
                        mLat.data.push(0);
                    }
                }
                mrs.logsync_lat.data.datasets.push(mLat);
                break;

            case "System":
                var mSys = valueui.utilx.objectClone(mTPL);
                switch (mTPL.label) {
                    case "CPU/Percent":
                        for (var j in m.points) {
                            if (m.points[j].count > 0) {
                                mSys.data.push(m.points[j].sum / m.points[j].count);
                            } else {
                                mSys.data.push(0);
                            }
                        }
                        mrs.system_cpu.data.datasets.push(mSys);
                        break;
                    case "Memory/Used":
                    case "Memory/Cached":
                        for (var j in m.points) {
                            if (m.points[j].count > 0) {
                                mSys.data.push((m.points[j].sum / m.points[j].count) / (1024 * 1024));
                            } else {
                                mSys.data.push(0);
                            }
                        }
                        mrs.system_mem.data.datasets.push(mSys);
                        break;
                    case "Net/Recv":
                    case "Net/Sent":
                    case "Disk/Read":
                    case "Disk/Write":
                        for (var j in m.points) {
                            if (j == 0) {
                                continue;
                            }
                            var s = m.points[j].sum - m.points[j - 1].sum;
                            if (s > 0) {
                                mSys.data.push(s / (1024 * 1024));
                            } else {
                                mSys.data.push(0);
                            }
                        }
                        mrs.system_io.data.datasets.push(mSys);
                        break;
                }
                break;
        }
    }

    for (var name in mrs) {
        if (update === true) {
            hooto_chart.RenderUpdate(mrs[name], "kvgo_instance_metric_" + name);
        } else {
            hooto_chart.RenderElement(mrs[name], "kvgo_instance_metric_" + name);
        }
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
                    id: "kvgo_instance_metric_service-qps",
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
    var elem = document.getElementById("kvgo_instance_metric_service_qps");
    if (!elem) {
        return ctx.callback("clean");
    }

    var req = "instance_name=" + kvgo.instanceActiveName + "&last_time_range=40&alignment_period=10";

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
