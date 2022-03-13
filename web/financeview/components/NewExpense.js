import React from 'react'

class NewExpense extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            date: dateFormat(),
            description: "",
            amount: 0,
            category: "",
            comment: ""
        }
        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleChange(event) {
        const target = event.target
        const value = target.value
        const name = target.name
        this.setState({
            [name]: value
        });
    }

    handleSubmit(event) {
        const query = {
            "query": `mutation NewExpense($date: String!, $desc: String!, $amt: Float!, $cats: [String!]!, $cmt: String) {createExpense(input:{
                date: $date,
                description: $desc,
                amount: $amt,
                categories: $cats,
                comment: $cmt
                }) {
                    Id
                    Date
                    Description
                    Amount
                    Categories {
                        Id
                        Name
                    }
                    Comment
                }  
            }`,
            variables: {
                date: this.state.date,
                desc: this.state.description,
                amt: parseFloat(this.state.amount),
                cats: [this.state.category],
                cmt: this.state.comment
            }
        }
        var myHeaders = new Headers();
        myHeaders.append("Content-Type", "application/json");

        var graphql = JSON.stringify(query)
        var requestOptions = {
            method: 'POST',
            headers: myHeaders,
            body: graphql,
            redirect: 'follow'
        };
        fetch("http://localhost:8080/query", requestOptions)
        .then(response => response.text())
        .then(result => console.log(result))
        .catch(error => console.log('error', error));
        // event.preventDefault();
    }

    render() {
        return (
            <form onSubmit={this.handleSubmit}>
                <label>
                    Date
                    <input
                        name="date"
                        type="text"
                        value={this.state.date}
                        onChange={this.handleChange} />
                </label>
                <label>
                    Description
                    <input
                        name="description"
                        type="text"
                        value={this.state.description}
                        onChange={this.handleChange} />
                </label>
                <label>
                    Amount
                    <input
                        name="amount"
                        type="number"
                        value={this.state.amount}
                        onChange={this.handleChange} />
                </label>
                <label>
                    Category
                    <input
                        name="category"
                        type="text"
                        value={this.state.category}
                        onChange={this.handleChange} />
                </label>
                <label>
                    Comment
                    <input
                        name="comment"
                        type="text"
                        value={this.state.comment}
                        onChange={this.handleChange} />
                </label>
                <input type="submit" value="Enter" />
            </form>
        )
    }
}

function dateFormat() {
    const dt = new Date()
    const year = dt.getFullYear().toString()
    const month = dt.getMonth().toString().padStart(2, '0')
    const day = dt.getDate().toString().padStart(2, '0')
    return month + "-" + day + "-" + year
}

export default NewExpense