var stat_table = false, statistics = {previous: [], current: []}, log_table = false;

$(document).ready(function () {
  $('#menu-containers').addClass('active');
  // $('#logs').css({maxHeight: '300px'});
});

var StatTableRow = React.createClass({
  render: function() {
    var stat = this.props.content.current,
        prev = this.props.content.previous,
        cpu_delta = 0, system_delta = 0, cpu_percent = 0;
    if (prev) {
      cpu_delta = stat.cpu_stats.cpu_usage.total_usage - prev.cpu_stats.cpu_usage.total_usage;
      system_delta = stat.cpu_stats.system_cpu_usage - prev.cpu_stats.system_cpu_usage;
    }
    if ((system_delta > 0) && (cpu_delta > 0)) {
      cpu_percent = 100.0 * cpu_delta / system_delta * stat.cpu_stats.cpu_usage.percpu_usage.length;
    }
    return (
        <tr key={this.props.index}>
          <td className="data-name">{stat.read.substring(5, 19).replace(/-/, '/').replace('T', ' ')}</td>
          <td className="data-name">{(cpu_percent+'').substring(0, 4)}%</td>
          <td className="data-name">{app.func.byteFormat(stat.memory_stats.usage)} / {app.func.byteFormat(stat.memory_stats.max_usage)} / {app.func.byteFormat(stat.memory_stats.limit)}</td>
          <td className="data-name">{((stat.memory_stats.usage * 100 / stat.memory_stats.limit)+'').substring(0, 4)}%</td>
          <td className="data-name">{app.func.byteFormat(stat.network.rx_bytes)} / {app.func.byteFormat(stat.network.tx_bytes)}</td>
          <td className="data-name">{stat.network.rx_packets} / {stat.network.tx_packets}</td>
        </tr>
    );
  }
});

var StatTable = React.createClass({
  getInitialState: function() {
    return {data: {previous: [], current: []}};
  },
  load: function(sender) {
    var id = $('#container-id').val();
    app.func.ajax({type: 'GET', url: '/api/container/stats/'+id, success: function (data) {
      statistics.previous = statistics.current;
      statistics.current = data;
      sender.setState({data: statistics});
      setTimeout(function () {stat_table && stat_table.setProps();}, 1000);
    }});
  },
  componentDidMount: function() {
    this.load(this);
  },
  componentWillReceiveProps: function() {
    this.load(this);
  },
  render: function() {
    var data = this.state.data, rows = [];
    $.map(data.current, function (current, index) {
      rows.push(<StatTableRow index={index} content={{previous: data.previous[index], current: current}} />)
    })
    return (
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th>Time</th>
              <th>CPU %</th>
              <th>MEM USAGE / MAX / LIMIT</th>
              <th>MEM %</th>
              <th>NET I/O</th>
              <th>NET I/O (packet)</th>
            </tr>
          </thead>
          <tbody>{rows}</tbody>
        </table>
    );
  }
});

stat_table = React.render(<StatTable />, document.getElementById('statistics'));


var LogTableRow = React.createClass({
  render: function() {
    var log = this.props.content;
    return (
        <tr>
          <td className="data-name no-wrap">{log.time}</td>
          <td className="data-name">{log.type}</td>
          <td className="data-name">{log.log}</td>
        </tr>
    );
  }
});

var LogTable = React.createClass({
  getInitialState: function() {
    return {data: []};
  },
  load: function(sender) {
    var id = $('#container-id').val();
    app.func.ajax({type: 'GET', url: '/api/container/logs/'+id, success: function (data) {
      var stream = [];
      $.map(data.stdout.split('\n'), function (record) {
        if (! record) return;
        stream.push({
          type: 'stdlog',
          key: record.substring(0, 30),
          time: record.substring(11, 19),
          log: record.substring(31)
        });
      });
      $.map(data.stderr.split('\n'), function (record) {
        if (! record) return;
        stream.push({
          type: 'stderr',
          key: record.substring(0, 30),
          time: record.substring(11, 19),
          log: record.substring(31)
        });
      });
      sender.setState({data: stream});
      setTimeout(function () {log_table && log_table.setProps();}, 1000);
    }});
  },
  componentDidMount: function() {
    this.load(this);
  },
  componentWillReceiveProps: function() {
    this.load(this);
  },
  render: function() {
    var rows = this.state.data.map(function(record, index) {
      if (! record.log) return;
      return <LogTableRow key={record.key} index={index} content={record} />
    });
    return (
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th>Time</th>
              <th>Type</th>
              <th>Log</th>
            </tr>
          </thead>
          <tbody>{rows}</tbody>
        </table>
    );
  }
});

log_table = React.render(<LogTable />, document.getElementById('logs'));
