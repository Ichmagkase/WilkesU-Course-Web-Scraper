import "./Card_Grid.css"
import Card from "../Card/Card.jsx"
import { useEffect, useState } from "react"

export default function Card_Grid({searchState, filterVisible}) {
  const [cards, setCards] = useState([]);


  // TODO: implement course validation by Filter
  const validateByFilter = (course) => {

  }

  // TODO: implement course validation by Search
  const validateBySearch = (course) => {
    const value = searchState.value.toLowerCase()
    return ((course.course_category ? String(course.course_category) : '').toLowerCase().includes(value) ||
            (course.instructor ? String(course.instructor) : '').toLowerCase().includes(value) ||
            (course.crn ? String(course.crn) : '').includes(value) ||
            (course.location ? String(course.location) : '').toLowerCase().includes(value) ||
            (course.title ? String(course.title) : '').toLowerCase().includes(value) ||
            (course.info ? String(course.info) : '').toLowerCase().includes(value))

    // return (String(course.course_category).toLowerCase().includes(value) ||
    //         String(course.instructor).toLowerCase().includes(value) ||
    //         String(course.crn).includes(value) ||
    //         String(course.location).toLowerCase().includes(value) ||
    //         String(course.title).toLowerCase().includes(value) ||
    //         String(course.info).toLowerCase().includes(value))
  }

  const GridWidth = {
    maxWidth: filterVisible ? "calc(100vw - 500px - 10px)" : "100vw"
  }

  useEffect( () => {
    const searchTerm = searchState;
    const tmpCards = [];

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

    // Get search parameters
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
        if (searchState.mode == "search") {
          for (let i = 0; i < courses.length; i++) {
            console.log(validateBySearch(courses[i]))
            if (validateBySearch(courses[i])) {
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
          }
        } else if (searchState.mode == "filter"){
          for (let i = 0; i < courses.length; i++) {
            if (validateByFilter(courses[i])) {
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
          }
        } else {
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
        }
        setCards(tmpCards)
      })
      .catch(error => console.error(error));
  }, [searchState]);

  return (
    <>
      <div className="grid_main" style={GridWidth}>
        {cards.map((props, index) => <Card key={index} {...props} />)}
      </div>
    </>
  )
}
