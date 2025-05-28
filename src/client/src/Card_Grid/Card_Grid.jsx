import "./Card_Grid.css"
import Card from "../Card/Card.jsx"
import { useEffect, useState } from "react"

export default function Card_Grid() {

  const [cards, setCards] = useState([]);

  useEffect( () => {
    /* Simulate Getting Data */
    const tmpCards = [];

    const filter = {
      semester: "F25",
      deliverymode: "",
      coursecategory: "",
      location: "",
      instructor: "",
      status: "",
    }

    const params = new URLSearchParams();

    for (const [key, value] of Object.entries(filter)) {
      if (value !== "") {
        params.append(key, value);
      }
    }

    const url = `http://localhost:8080/filter?${params.toString()}`;

    let courses;
    fetch(url)
      .then(response => response.json())
      .then(data => {
        courses = data
        for (let i = 0; i < courses.length; i++) {
          tmpCards.push({
            header: `${courses[i].course_category + " " + courses[i].course_id}`,
            title: `${courses[i].title}`,
            credits: `${courses[i].credits}`,
            extra_info: `${courses[i].info}`,
            time: `${courses[i].start_time} - ${courses[i].end_time + courses[i].end_time_ampm}`,
            crn: `${courses[i].crn}`,
            students: `${courses[i].students}`,
            limit: `${courses[i].limit}`
          });
        }
        setCards(tmpCards)
      })
      .catch(error => console.error(error));
  }, []);



  return (
    <>
      <div className="grid_main">
        {cards.map((props, index) => <Card key={index} {...props} />)}
      </div>
    </>
  )
}
