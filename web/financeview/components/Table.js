import React from 'react'
class Table extends React.Component {
    constructor(props) {
        super(props)
        var cols = ['test', 'column1']
        var data = [['a', 'b'], ['1',2]]
        const exps = props.data
        console.log("printing exps")
        console.log(exps)
        // const elm = exps[0];
        // console.log(elm)
        // for (let key in Object.keys(elm)) {
        //     cols.push(key)
        // }
        // for (let elm in exps) {
        //     row = []
        //     for (let key in cols) {
        //         row.push(elm[key])
        //     }
        //     data.push(row)
        // }
        this.state = {
            cols: cols,
            data: data
        }
    }

    render() {
        return (
            <table>
                <thead>
                    <tr>
                        {this.state.cols.map(name => <th key={name}>{name}</th>)}
                    </tr>
                </thead>
                <tbody>
                    {this.state.data.map(row => <tr key={row}>{row.map(data => <td key={data}>{data}</td>)}</tr>)}
                </tbody>
            </table>
        )
    }
}

export default Table;