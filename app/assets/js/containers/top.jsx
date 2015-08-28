$(document).ready(function () {
  $('#menu-containers').addClass('active');
});

var TableRow = React.createClass({
  render: function() {
    var process = this.props.content;
    return (
        <tr key={this.props.index}>
          <td className="data-name">{process[0]}</td>
          <td className="data-name">{process[1]}</td>
          <td className="data-name">{process[2]}</td>
          <td className="data-name">{process[3]}</td>
          <td className="data-name">{process[4]}</td>
          <td className="data-name">{process[5]}</td>
          <td className="data-name">{process[6]}</td>
          <td className="data-name">{process[7]}</td>
          <td className="data-name">{process[8]}</td>
          <td className="data-name">{process[9]}</td>
          <td className="data-name">{process[10]}</td>
        </tr>
    );
  }
});

var Table = React.createClass({
  getInitialState: function() {
    return {data: {}};
  },
  componentDidMount: function() {
    var self = this, id = $('#container-id').val();
    app.func.ajax({type: 'GET', url: '/api/container/top/'+id, success: function (data) {
      self.setState({data: data});
    }});
  },
  render: function() {
    if ((! this.state.data) || (! this.state.data.Processes)) {
      return (
        <table className="table table-striped table-hover">
          <thead></thead>
          <tbody><tr><th>There is no data.</th></tr></tbody>
        </table>
      );
    }
    var titles = this.state.data.Titles;
    var rows = this.state.data.Processes.map(function(record, index) {
      return <TableRow index={index} content={record} />
    });
    return (
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th>{titles[0]}</th><th>{titles[1]}</th><th>{titles[2]}</th><th>{titles[3]}</th>
              <th>{titles[4]}</th><th>{titles[5]}</th><th>{titles[6]}</th><th>{titles[7]}</th>
              <th>{titles[8]}</th><th>{titles[9]}</th><th>{titles[10]}</th>
            </tr>
          </thead>
          <tbody>{rows}</tbody>
        </table>
    );
  }
});

React.render(<Table />, document.getElementById('data'));
