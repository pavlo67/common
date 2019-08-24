import React, {Component} from 'react'
import 'bootstrap/dist/css/bootstrap.css'
import {Menu, Workspace} from './menu'

class App extends Component {
    // state = {
    //     reverted: false
    // };

    render() {
        // console.log('---', 1)

        return (
            <div className="App">
                <Menu />
                <Workspace />
            </div>
        );
    }

    // revert = () => {
    //     this.setState({
    //         reverted: !this.state.reverted
    //     })
    // }
}

export default App


// import ArticleList from './old/ArticleList/index'
// import articles from './old/fixtures'
// <div className="container">
//     <div className="jumbotron">
//     <h1 className="display-3">
//     App name
// <button className="btn" onClick = {this.revert}>Revert</button>
// </h1>
// </div>
// <ArticleList articles={this.state.reverted ? articles.slice().reverse() : articles}/>
// </div>
