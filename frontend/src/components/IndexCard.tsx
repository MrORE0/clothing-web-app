import "../styles/index-card.css"
import modelsImage from "../assets/Model_Colorized.jpeg"

export const IndexCard = () => {
    return (
        <>
        <div className="card text-white rounded-0 border-0" style={{"maxWidth": "1600px"}}>
            <img src={modelsImage} className="" alt="Models Image"/>
            <div className="card-img-overlay rounded-0 d-flex flex-column index-card-text" style={{"backgroundColor": "rgba(0, 0, 0, 0.3)"}}>
                <p className="card-text fs-60 fw-bolder mt-4 font-style-fjallaone">
                    A way for you to dress like <span className="fs-70 text-decoration-underline">you</span>.
                </p>
                <p className="card-text fs-70 fw-bolder mt-4 font-style-fjallaone text-end">Choose <span className="fs-80 text-decoration-underline">your</span> fit.</p>
            </div>
        </div>
        </>
    )
}