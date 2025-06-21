import Header from "./Header/Header.jsx"
import Card_Grid from "./Card_Grid/Card_Grid.jsx"
import Footer from "./Footer/Footer.jsx"
import { useState } from 'react';

export default function App() {

  const [searchState, setSearchTerm] = useState({})
  const [filterVisible, setFilterVisible] = useState(false)
  console.log("updated state")

  return (
    <>
      <Header setSearchTerm={setSearchTerm} filterVisible={filterVisible} setFilterVisible={setFilterVisible}/>
      <Card_Grid searchState={searchState} filterVisible={filterVisible}/>
      <Footer/>
    </>
  )
}
