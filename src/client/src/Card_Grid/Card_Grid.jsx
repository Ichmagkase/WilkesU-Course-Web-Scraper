import "./Card_Grid.css"
import Card from "../Card/Card.jsx"
import { useEffect, useState } from "react"
function Card_Grid() {

  const [cards, setCards] = useState([]);

  useEffect( () => {
    /* Simulate Getting Data */
    const tmpCards = [];

    fetch("http://localhost:8080/filter?semester=F25&instructor=kapolka")
      .then(response => response.json())
      .then(data => console.log(data))
      .catch(error => console.error(error));

    for (let i = 0; i < 10; i++) {
      tmpCards.push({
        header: `Header ${i}`,
        title: `Title ${i}`,
        credits: `${i}`,
        extra_info: `Extra Info ${i}`,
        time: `Time ${i}`,
        crn: `${i}0000`,
        students: `${i}`,
        limit: `${i+1}`
      });
      let coursedata
    }
    setCards(tmpCards)
  }, []);

  return (
    <>
      <div className="grid_main">
        {cards.map((props, index) => <Card key={index} {...props} />)}
      </div>
    </>
  )
}

export default Card_Grid
