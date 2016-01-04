$(document).ready(function () {
  $('#menu-containers').addClass('active');
});

var TableRow = React.createClass({
  render: function() {
    var change = this.props.content,
        kind = '';
    switch (change.Kind) {
      case 0: kind = 'Modify'; break;
      case 1: kind = 'Add'; break;
      case 2: kind = 'Delete'; break;
    }
    return (
        <tr key={this.props.index}>
          <td className="data-name">{kind}</td>
          <td className="data-name">{change.Path}</td>
        </tr>
    );
  }
});

var Table = React.createClass({
  getInitialState: function() {
    return {data: []};
  },
  componentDidMount: function() {
    var self = this,
        id = $('#container-id').val(),
        client = $('#client-id').val();
    client = client ? '?client='+client : '';
    app.func.ajax({type: 'GET', url: '/api/container/changes/'+id+client, success: function (data) {
      self.setState({data: data});
    }});
  },
  render: function() {
    var rows;
    if (this.state.data.length > 0) {
      rows = this.state.data.map(function(record, index) {
        return <TableRow index={index} content={record} />
      });
    } else {
      rows = <TableRow index={0} content={{Kind: -1, Path: 'There is no changed file.'}} />
    }
    return (
        <table className="table table-striped table-hover">
          <thead>
            <tr>
              <th></th>
              <th>Path</th>
            </tr>
          </thead>
          <tbody>{rows}</tbody>
        </table>
    );
  }
});

ReactDOM.render(<Table />, document.getElementById('data'));
