import Link from 'next/Link'
import Head from 'next/head'
import Table from '../../components/Table'
import styles from '../../styles.module.css'
import React from 'react'

class Expenses extends React.Component {
    constructor(props) {
        super(props)
        this.state = {
            cols: [],
            data: []
        }
    }

    componentDidMount() {
        listExpenses().then(response => {
            let keys = []
            for (const key in response.expenses[0]) {
                keys.push(key)
            }
            let data = []
            for (const elm in response.expenses) {
                let row = []
                for (var i=0; i < keys.length; i++) {
                    var col = keys[i]
                    if (col === "Categories") {
                        var cats = response.expenses[elm][col]
                        var val = ""
                        for (var j=0; j < cats.length; j++) {
                            if (j != 0) {
                                val += ", "
                            }
                            val += cats[j].Name
                        }
                        row.push(val)
                    } else {
                        row.push(response.expenses[elm][col])
                    }
                }
                data.push(row)
            }
            this.setState({
                cols: keys,
                data: data
            })
        });
    }
    render() {
        return (
            <>
            <Head>
                <title>FinanceView</title>
            </Head>
            <h1 className={styles.header}>Expenses Page!</h1>
            <h2>
                <Link href="/"><a>Home</a></Link>
            </h2>
            <br></br>
            <Table cols={this.state.cols} data={this.state.data} />
        </>
        )
    }
}


export async function listExpenses() {
    const query = {"query": ` { expenses {
        Id
        Date
        Amount
        Description
        Categories {
            Id
            Name
        }
        Comment
    }}`}
    const res = await fetch(
        'http://localhost:8080/query',
        {
            method: "POST",
            body: JSON.stringify(query),
            headers: new Headers({'content-type': 'application/json'})
        }
    )
    const json = await res.json()
    const data = json.data
    return data
}

export default Expenses