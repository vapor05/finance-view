import Link from 'next/Link'
import Head from 'next/head'
import styles from '../styles.module.css'

export default function HomePage() {
    return (
        <>
            <Head>
                <title>FinanceView</title>
            </Head>
            <h1 className={styles.header}>Welcome to FinanceView!</h1>
            
            <Link href="/expenses"><a>See expenses</a></Link>
        </>
    )
}

