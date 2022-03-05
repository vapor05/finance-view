import Link from 'next/Link'
import Head from 'next/head'
import styles from '../../styles.module.css'

export default function Expenses() {
    const data = getData()
    console.log(data)
    return (
        <>
            <Head>
                <title>FinanceView</title>
            </Head>
            <h1 className={styles.header}>Expenses Page!</h1>
            <h2>
                <Link href="/"><a>Home</a></Link>
            </h2>
            
        </>
    )
}

export async function getData() {
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
    console.log(json)
    return {
        props: { json }
    }
}