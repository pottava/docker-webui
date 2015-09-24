var table = false,
    filters = {client: app.func.query('c', -1), label: 0},
    allcontainers = [],
    statistics = {
      previous: [],
      current: []
    },
    storedCPUs = [],
    storedMems = [],
    refreshWindow = app.storage.get('refresh-statistics', 2),
    clients = 0,
    countup = 0;

$(document).ready(function () {
  $('#menu-statistics').addClass('active');
  clients = parseInt($('#number-of-clients').val(), 10);

  setRefreshWindow(refreshWindow);
  $('#refresh-window a').click(function(e) {
    setRefreshWindow(parseInt($(this).attr('href').substring(1), 10));
    app.func.stop(e);
  });

  // make filters
  app.func.ajax({url: '/api/containers', data: {status: 3}, success: function (data) {
    allcontainers = data;
    var clients = [],
        labels = [];

    var temp = {clients: {}, labels: {}};
    $.map(allcontainers, function (candidate) {
      temp.clients[''+candidate.client.id] = candidate.client.endpoint.replace(/^.*:\/\//, '').replace(/:.*$/, '');

      $.map(candidate.containers, function (container) {
        if (! container.labels) return;
        $.map(container.labels, function (value, key) {
          temp.labels[key] = value;
        });
      });
    });
    $.map(temp.clients, function (value, key) {
      clients.push({key: key, value, value});
    });
    $.map(temp.labels, function (value, key) {
      labels.push({key: key, value, value});
    });
    clients.sort(function (a, b) {
      if (a.value < b.value) return -1;
      if (a.value > b.value) return 1;
      return 0;
    });
    labels.sort(function (a, b) {
      if (a.key < b.key) return -1;
      if (a.key > b.key) return 1;
      if (a.value < b.value) return -1;
      if (a.value > b.value) return 1;
      return 0;
    });
    _setClientOption(clients);
    _setLabelFilter(labels);

    refreshStats();
  }});
});

function setRefreshWindow(value) {
  var a = $('#refresh-window a[href="#'+value+'"]'),
      group = a.closest('.btn-group').removeClass('open');
  refreshWindow = value;
  app.storage.set('refresh-statistics', value);
  group.find('.caption').text('refresh / '+a.text()).blur();
}

function refreshStats() {
  if (refreshWindow && (refreshWindow > 0)) {
    table && table.setProps();
  }
  setTimeout(refreshStats, Math.max(1, refreshWindow) * 1000);
}

function _find(arr, key, def) {
  var result = def ? def : '';
  if (! arr) {
    return result;
  }
  $.map(arr, function (value) {
    if (value.split('=')[0] == key) {
      result = value.split('=')[1];
    }
  });
  return result;
}

function _setClientOption(clients) {
  var options = $('.client-filters').hide(),
      caption = false,
      count = 0;
  options.find('ul.dropdown-menu').html('');
  $.map(clients, function (client) {
    if ((! caption) && (filters.client == client.key)) caption = client.value;
    count++;
  });
  caption && options.find('.caption').text(caption);
  if (count <= 1) {
    $('.line-charts').css({width:'70%'});
    $('.pie-charts').show();
    return;
  }
  var html = '<li><a href="#-1">All</a></li>';
  $.map(clients, function (client) {
    html += '<li><a href="#'+client.key+'">'+client.value+'</a></li>';
  });
  options.find('ul.dropdown-menu').html(html);
  options.fadeIn();

  $('#client-filter .dropdown-menu a').click(function(e) {
    var a = $(this),
        group = a.closest('.btn-group').removeClass('open');
    group.find('.caption').text(a.text()).blur();
    filters.client = a.attr('href').substring(1);

    if (filters.client != -1) {
      $('.line-charts').css({width:'70%'});
      $('.pie-charts').show();
    } else {
      $('.pie-charts').hide();
      $('.line-charts').css({width:'95%'});
    }
    app.func.stop(e);
  });
}

function _setLabelFilter(labels) {
  var options = $('.label-filters').hide(),
      caption = false,
      count = 0;
  options.find('ul.dropdown-menu').html('');
  $.map(labels, function (label) {
    if ((! caption) && (filters.label == label.key)) caption = label.value;
    count++;
  });
  if ((! caption) && (filters.label == -1)) caption = 'Not Labeled';
  caption && options.find('.caption').text(caption);
  if (count <= 1) {
    return;
  }
  var html = '<li><a href="#0">All</a></li>',
      group = '';
  html += '<li><a href="#-1">Not Labeled</a></li>';
  $.map(labels, function (label) {
    if (group != label.key) {
      html += '<li class="dropdown-header">'+label.key+'</li>';
      group = label.key;
    }
    var value = label.value;
    if (value.length > 20) {
      value = value.substring(0, 20) + '..';
    }
    value = '&nbsp;&nbsp;&nbsp;&nbsp;'+value;
    html += '<li><a href="#'+app.func.hash(label.key+'->'+label.value)+'">'+value+'</a></li>';
  });
  options.find('ul.dropdown-menu').html(html);
  options.fadeIn();

  $('#label-filter .dropdown-menu a').click(function(e) {
    var a = $(this),
        group = a.closest('.btn-group').removeClass('open');
    group.find('.caption').text(a.text().trim()).blur();
    filters.label = a.attr('href').substring(1);
    app.func.stop(e);
  });
}

var lineChartCPU = nv.models.lineChart()
              .x(function(d) { return d[0] })
              .y(function(d) { return d[1]/100 })
              .color(d3.scale.category10().range())
              .showLegend(false)
              .useInteractiveGuideline(true);
lineChartCPU.xAxis.tickFormat(function(d) {return d3.time.format('%X')(new Date(d))});
lineChartCPU.yAxis.tickFormat(d3.format(',.1%'));
nv.utils.windowResize(lineChartCPU.update);

var lineChartMem = nv.models.lineChart()
              .x(function(d) { return d[0] })
              .y(function(d) { return d[1]/100 })
              .color(d3.scale.category10().range())
              .showLegend(false)
              .useInteractiveGuideline(true);
lineChartMem.xAxis.tickFormat(function(d) {return d3.time.format('%X')(new Date(d))});
lineChartMem.yAxis.tickFormat(d3.format(',.1%'));
nv.utils.windowResize(lineChartMem.update);

var pieChartCPU = nv.models.pieChart()
            .x(function(d) { return d.label })
            .y(function(d) { return d.value })
            .color(d3.scale.category10().range())
            .showLegend(false)
            .showLabels(true)
            .labelThreshold(.05)
            .labelType("percent")
            .donut(true)
            .donutRatio(0.35);

var pieChartMem = nv.models.pieChart()
            .x(function(d) { return d.label })
            .y(function(d) { return d.value })
            .color(d3.scale.category10().range())
            .showLegend(false)
            .showLabels(true)
            .labelThreshold(.05)
            .labelType("percent")
            .donut(true)
            .donutRatio(0.35);

function _setStoredData(arr, index, key, values) {
  var found = false;
  $.map(arr, function (record) {
    if (record.key == key) found = record;
  });
  if (found) {
    found.indexes.push(index);
    found.values.push(values);
    return;
  }
  arr.push({key: key, indexes: [index], values: [values]});
}

function _shiftStoredData(arr, index) {
  $.map(arr, function (record) {
    if (record.indexes[0] <= (index - 20)) {
      record.indexes.shift();
      record.values.shift();
    }
  });
}

function _spliceStoredData(arr, names) {
  for(var i = arr.length - 1; i >= 0; i--) {
    if (! names[arr[i].key]) {
      arr.splice(i, 1);
    }
  }
}

var TableRow = React.createClass({
  render: function() {
    var name = this.props.content.name,
        stat = this.props.content.current && this.props.content.current,
        prev = this.props.content.previous && this.props.content.previous[0],
        cpu_delta = 0, system_delta = 0, cpu_percent = 0;
    if (prev && stat && stat.cpu_stats) {
      cpu_delta = stat.cpu_stats.cpu_usage.total_usage - prev.cpu_stats.cpu_usage.total_usage;
      system_delta = stat.cpu_stats.system_cpu_usage - prev.cpu_stats.system_cpu_usage;
    }
    if ((system_delta > 0) && (cpu_delta > 0)) {
      cpu_percent = 100.0 * cpu_delta / system_delta * stat.cpu_stats.cpu_usage.percpu_usage.length;
    }
    var time = '',
        mem = {usage: '-', max: '-', limit: '-', percent: 0},
        network = {in: '', out: '', inPacket: '', outPacket: ''};
    if (stat && stat.read) {
      time = stat.read.substring(5, 19).replace(/-/, '/').replace('T', ' ');
      mem = {
        usage: app.func.byteFormat(stat.memory_stats.usage),
        max: app.func.byteFormat(stat.memory_stats.max_usage),
        percent: stat.memory_stats.usage * 100 / stat.memory_stats.limit
      };
      network = {
        in: app.func.byteFormat(stat.network.rx_bytes),
        out: app.func.byteFormat(stat.network.tx_bytes)
      };
    }
    return (
        <tr key={this.props.index}>
          <td className="data-name">{name.substring(1).replace(',/', ', ')}</td>
          <td className="data-name">{time}</td>
          <td className="data-name">{(cpu_percent+'').substring(0, 4)}%</td>
          <td className="data-name">{mem.usage} / {mem.max}</td>
          <td className="data-name">{(mem.percent+'').substring(0, 4)}%</td>
          <td className="data-name">{network.in} / {network.out}</td>
        </tr>
    );
  }
});

var Table = React.createClass({
  getInitialState: function() {
    return {data: {stats: [], multiple: false}};
  },
  load: function(sender) {
    var self = this, conditions = {};
    if (filters.client != -1) conditions.client = filters.client;

    app.func.ajax({url: '/api/statistics', data: conditions, success: function (data) {
      countup++;

      data.sort(function (a, b) {
        if (a.client.endpoint < b.client.endpoint)
          return -1;
        if (a.client.endpoint > b.client.endpoint)
          return 1;
        return 0;
      });

      // filter by labels
      var allowedIds = {};
      $.map(allcontainers, function (candidate) {
        $.map(candidate.containers, function (container) {
          if ((filters.label == 0) || ((filters.label == -1) && (! container.labels))) {
            allowedIds[container.id] = true;
          } else {
            if (! container.labels) return;
            var match = false;
            $.map(container.labels, function (value, key) {
              match |= (filters.label == app.func.hash(key+'->'+value));
            });
            if (match) allowedIds[container.id] = true;
          }
        });
      });

      // backup data for previous calculation
      if (data.error) {
        statistics.previous = true;
        statistics.current = data;
      } else {
        statistics.previous = statistics.current;
        statistics.current = data;
      }

      // retrieve stats data
      var stats = [];
      $.map(data, function (host) {
        $.map(host.stats, function (nest, id) {
          if (! allowedIds[id]) return;

          $.map(nest, function (record, key) {
            stats.push({
              client: host.client,
              id: id, key: key,
              stat: record && record[0]
            })
          });
        });
      });
      stats.sort(function (a, b) {
        var an = a.key + a.client.endpoint,
            bn = b.key + b.client.endpoint;
        if (an < bn) return -1;
        if (an > bn) return  1;
        return 0;
      });

      // update stats-table
      sender.setState({data: {stats: stats, multiple: (data.length > 1)}});

      // change data format for charts
      var names = {},
          pie = {CPU: [], Mem: []},
          multiple = (data.length > 1);
      $.map(stats, function (record, index) {
        var client = record.client,
            name = record.key.substring(1).replace(',/', ', ') + _endpoint(multiple, client.endpoint),
            stat = record.stat,
            time = new Date(stat.read.substring(0, 19)+'Z').getTime(),
            prev = _findPrivious(client.endpoint, record.id, record.key),
            cpu_delta = 0, system_delta = 0, cpu_percent = 0, mem_percent = 0;
        if (prev && (prev.length > 0)) {
          prev = prev[0];
        }
        if (prev && stat && stat.cpu_stats) {
          cpu_delta = stat.cpu_stats.cpu_usage.total_usage - prev.cpu_stats.cpu_usage.total_usage;
          system_delta = stat.cpu_stats.system_cpu_usage - prev.cpu_stats.system_cpu_usage;
        }
        if ((system_delta > 0) && (cpu_delta > 0)) {
          cpu_percent = 100.0 * cpu_delta / system_delta * stat.cpu_stats.cpu_usage.percpu_usage.length;
        }
        if (stat) mem_percent = stat.memory_stats.usage * 100 / stat.memory_stats.limit;

        _setStoredData(storedCPUs, countup, name, [time, cpu_percent]);
        pie.CPU.push({label: name, value: cpu_percent});
        _setStoredData(storedMems, countup, name, [time, mem_percent]);
        pie.Mem.push({label: name, value: mem_percent});
        names[name] = true;
      });
      _shiftStoredData(storedCPUs, countup);
      _shiftStoredData(storedMems, countup);
      _spliceStoredData(storedCPUs, names);
      _spliceStoredData(storedMems, names);

      // draw line-charts
      nv.addGraph(function() {
        d3.select('#chart-cpu svg.line-charts').datum(storedCPUs).call(lineChartCPU);
        var max = d3.max(storedCPUs, function(d) {
          return d3.max(d.values, function(value) {
            return (value[1] + 0.1) / 100;
          });
        });
        return lineChartCPU.yDomain([0, max]);
      });
      if ((clients = 1) || (filters.client != -1)) {
        var remain = 100;
        $.map(pie.CPU, function (record) {
          remain -= record.value;
        });
        pie.CPU.push({label: '-', value: remain});
        d3.select("#chart-cpu svg.pie-charts").datum(pie.CPU).transition().duration(350).call(pieChartCPU);
      }
      nv.addGraph(function() {
        d3.select('#chart-mem svg.line-charts').datum(storedMems).call(lineChartMem);
        var max = d3.max(storedMems, function(d) {
          return (d.values[d.values.length-1][1] + 0.5) / 100;
        });
        return lineChartMem.yDomain([0, max]);
      });
      if ((clients = 1) || (filters.client != -1)) {
        var remain = 100;
        $.map(pie.Mem, function (record) {
          remain -= record.value;
        });
        pie.Mem.push({label: '-', value: remain});
        d3.select("#chart-mem svg.pie-charts").datum(pie.Mem).transition().duration(350).call(pieChartMem);
      }
    }});
  },
  componentDidMount: function() {
    this.load(this);
  },
  componentWillReceiveProps: function() {
    this.load(this);
  },
  render: function() {
    var multiple = this.state.data.multiple,
        data = this.state.data.stats,
        rows = [];
    $.map(data, function (record, index) {
      var client = record.client,
          name = record.key + _endpoint(multiple, client.endpoint);
      rows.push(<TableRow key={record.key+'@'+client.id} index={record.key+'@'+index} content={{
        name: name,
        current: record.stat,
        previous: _findPrivious(client.endpoint, record.id, record.key)
      }} />)
    });
    return (
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th>ID</th>
              <th>Time</th>
              <th>CPU %</th>
              <th>MEM USAGE / MAX</th>
              <th>MEM %</th>
              <th>NET I/O</th>
            </tr>
          </thead>
          <tbody>{rows}</tbody>
        </table>
    );
  }
});

function _endpoint(multiple, endpoint) {
  return (multiple ? ' @'+endpoint.replace(/^.*:\/\//, '').replace(/:.*$/, '') : '');
}

function _findPrivious(endpoint, id, key) {
  var result = false;
  if (statistics.previous) {
    $.map(statistics.previous, function (host) {
      if (host.client.endpoint == endpoint)
        result = host.stats[id][key];
    });
  }
  return result;
}

table = React.render(<Table />, document.getElementById('data'));
