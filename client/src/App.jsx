import { useState } from 'react'
import Card from "./Card/Card.jsx"

function App() {
  return (
    <>
      <Card/>
      <Card 
        header="F2F - PHA 425 (A) Kieck D"
        title="Pharmacotherapeutics III"
        credits="3.00"
        time="M 2 - 3:50PM SLC B05; TRF 2 - 3:50PM SLC 105"
        crn="31146"
        students="55"
        limit="105"
      />
    </>
  )
}

export default App
