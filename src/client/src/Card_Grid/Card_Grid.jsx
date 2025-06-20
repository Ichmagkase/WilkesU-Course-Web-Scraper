import "./Card_Grid.css"
import Card from "../Card/Card.jsx"
import { useEffect, useState } from "react"

export default function Card_Grid({searchState}) {
  const [cards, setCards] = useState([]);

  useEffect( () => {
    const searchTerm = searchState;
    const tmpCards = [];

    console.log("load");

    const serverFilter = {
      semester: "F25",

      // The following filter options may also be enabled for
      // serverside filtering (via a DB call)

      // deliverymode: "",
      // coursecategory: "",
      // location: "",
      // instructor: "",
      // status: "",
    }

    const params = new URLSearchParams();

    for (const [key, value] of Object.entries(serverFilter)) {
      if (value !== "") {
        params.append(key, value);
      }
    }

    const url = `http://localhost:8080/filter?${params.toString()}`;
    let courses;
    fetch(url)
      .then(response => response.json())
      .then(data => {
        courses = data;
        for (let i = 0; i < courses.length; i++) {
          tmpCards.push({
            header: `${courses[i].course_category + " " + courses[i].course_id}`,
            instructor: `${courses[i].instructor}`,
            section: `${courses[i].course_section}`,
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
  }, [searchState]);

  return (
    <>
      <div className="grid_main">
        {cards.map((props, index) => <Card key={index} {...props} />)}
      </div>
    </>
  )
}
