import "./Card_Grid.css"
import Card from "../Card/Card.jsx"
import { useEffect, useState } from "react"

export default function Card_Grid({searchState, filterVisible}) {
  const [cards, setCards] = useState([]);

  // TODO: implement course validation by Filter
  const validateByFilter = (course) => {

  }

  const validateBySearch = (course) => {
    const value = searchState.value.toLowerCase()

    return ((course.course_category ? String(course.course_category) + " " + String(course.course_id)  : '').toLowerCase().includes(value) ||
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

  // Simple debouncing implementation by nishant-666:
  // https://codesandbox.io/p/sandbox/react-debouncing-k5qdlv?file=%2Fsrc%2FApp.js
  useEffect( () => {
    const initiateSearch = setTimeout(() => {
      const searchTerm = searchState;
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
          const courses = data;
          let tmpCards = [];
          const addCard = (course) => {
            const category = course.course_category;
            const courseid = course.course_id ? course.course_id : course.course_id = "";
            const extra_info = course.info ? course.info : course.info = "";
            let time = course.start_time ? course.start_time : course.start_time = "TBA";
            time = String(course.start_time).includes("TBA") ? "Time: TBA" : `${course.start_time} - ${course.end_time + course.end_time_ampm}`;
            tmpCards.push({
              header: `${course.course_category} ${course.course_id}`,
              instructor: `${course.instructor}`,
              section: `${course.course_section}`,
              title: `${course.title}`,
              credits: course.credits,
              extra_info: `${course.info}`,
              time: time,
              crn: course.crn,
              students: course.students,
              limit: course.limit
            });
          };

          if (searchState?.mode === "search") {
            courses.filter(validateBySearch).forEach(addCard);
          } else if (searchState?.mode === "filter") {
            courses.filter(validateByFilter).forEach(addCard);
          } else {
            courses.forEach(addCard);
          }

          setCards(tmpCards);

        })
        .catch(error => console.error(error));
    }, 500);
    return () => clearTimeout(initiateSearch);
  }, [searchState]);

  return (
    <>
      <div className="grid_main" style={GridWidth}>
        {cards.map((props, index) => <Card key={index} {...props} />)}
      </div>
    </>
  )
}
