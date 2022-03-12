import React from 'react'
import styles from './table.module.css'

class Table extends React.Component {
    constructor(props) {
        super(props)
        console.log("printing props")
        console.log(props)
        this.state = {
            cols: props.cols,
            data: props.data
        }
    }

    render() {
        return (
            <table className={styles.table}>
                <thead>
                    <tr>
                        {this.props.cols.map(name => <th className={styles.th} key={name}>{name}</th>)}
                    </tr>
                </thead>
                <tbody>
                    {this.props.data.map((row, index) => <tr key={index}>{row.map((data, index) => <td className={styles.td} key={index}>{data}</td>)}</tr>)}
                </tbody>
            </table>
        )
    }
}

export default Table;