import Header from "./Header/Header.jsx"
import Card_Grid from "./Card_Grid/Card_Grid.jsx"
import Footer from "./Footer/Footer.jsx"
import { useState } from 'react';

export default function App() {

  const [clientSearch, setSearchTerm] = useState({})
  console.log("App: ")
  console.log(clientSearch);

  return (
    <>
      <Header setFilter={setSearchTerm}/>
      <Card_Grid searchState={clientSearch} className="card-grid-wrapper"/>
      <Footer/>
    </>
  )
}
