package controllers

import (
	"github.com/jgraham909/bloggo/app/models"
	"github.com/jgraham909/revmgo"
	"github.com/robfig/revel"
	"github.com/russross/blackfriday"
	"html/template"
	"labix.org/v2/mgo/bson"
	"time"
)

type Application struct {
	*revel.Controller
	revmgo.MongoController
	ActiveUser *models.User
}

// Responsible for doing any necessary setup for each web request.
func (c *Application) Setup() revel.Result {
	// If there is an active user session load the User data for this user.
	if email, ok := c.Session["user"]; ok {
		c.ActiveUser = models.GetUserByEmail(c.MongoSession, email)
		c.RenderArgs["ActiveUser"] = c.ActiveUser
	}

	dummyContent()
	return nil
}

func (c Application) Index() revel.Result {
	return c.Redirect(Blog.Index)
}

func (c Application) UserAuthenticated() bool {
	if _, ok := c.Session["user"]; ok {
		return true
	}
	return false
}

func (c Application) Preview(text string) revel.Result {
	preview := template.HTML(string(blackfriday.MarkdownBasic([]byte(text))))
	return c.Render(preview)

}

func (c Application) NavLeft() revel.Result {
	UserCreate := c.ActiveUser.CanBeCreatedBy(c.MongoSession, c.ActiveUser)

	article := new(models.Article)
	ArticleCreate := article.CanBeCreatedBy(c.MongoSession, c.ActiveUser)
	return c.Render(UserCreate, ArticleCreate)
}

func (c Application) NavRight() revel.Result {
	UserUpdate, UserLogout, UserLogin := false, false, true
	if c.ActiveUser != nil {
		UserUpdate = c.ActiveUser.CanBeUpdatedBy(c.MongoSession, c.ActiveUser)
		UserLogout = c.ActiveUser.Id.Valid()
		UserLogin = !c.ActiveUser.Id.Valid()
	}

	return c.Render(UserUpdate, UserLogout, UserLogin)
}

// Load dummy content if the jane account doesn't exist.
func dummyContent() {
	jane := models.GetUserByObjectId(revmgo.Session, bson.ObjectIdHex("51e9aa4049a1b716bb000003"))
	if jane.Id.Hex() == "" {
		jane = new(models.User)
		jane.Id = bson.ObjectIdHex("51e9aa4049a1b716bb000003")
		jane.Firstname = "Jane"
		jane.Lastname = "Doe"
		jane.Email = "jane@example.com"
		p := models.Password{"12341234", "12341234"}
		jane.Save(revmgo.Session, p)

		john := new(models.User)
		john.Id = bson.ObjectIdHex("51e9aa2d49a1b716bb000002")
		john.Firstname = "John"
		john.Lastname = "Doe"
		john.Email = "john@example.com"
		john.Save(revmgo.Session, p)

		a := new(models.Article)
		a.Id = bson.ObjectIdHex("51e9ad9749a1b71843000001")
		a.Author_id = bson.ObjectIdHex("51e9aa4049a1b716bb000003")
		a.Published = true
		a.Posted = time.Date(2013, time.July, 19, 12, 32, 00, 0, time.UTC)
		a.Title = "Moby Dick"
		a.Body = "<p>Queequeg was a native of Rokovoko, an island far away to the West and South. It is not down in any map; true places never are.</p>\r\n\r\n<p>When a new-hatched savage running wild about his native woodlands in a grass clout, followed by the nibbling goats, as if he were a green sapling; even then, in Queequeg's ambitious soul, lurked a strong desire to see something more of Christendom than a specimen whaler or two. His father was a High Chief, a King; his uncle a High Priest; and on the maternal side he boasted aunts who were the wives of unconquerable warriors. There was excellent blood in his veins&mdash;royal stuff; though sadly vitiated, I fear, by the cannibal propensity he nourished in his untutored youth.</p>\r\n\r\n<p>A Sag Harbor ship visited his father's bay, and Queequeg sought a passage to Christian lands. But the ship, having her full complement of seamen, spurned his suit; and not all the King his father's influence could prevail. But Queequeg vowed a vow. Alone in his canoe, he paddled off to a distant strait, which he knew the ship must pass through when she quitted the island. On one side was a coral reef; on the other a low tongue of land, covered with mangrove thickets that grew out into the water. Hiding his canoe, still afloat, among these thickets, with its prow seaward, he sat down in the stern, paddle low in hand; and when the ship was gliding by, like a flash he darted out; gained her side; with one backward dash of his foot capsized and sank his canoe; climbed up the chains; and throwing himself at full length upon the deck, grappled a ring-bolt there, and swore not to let it go, though hacked in pieces.</p>\r\n\r\n<p>In vain the captain threatened to throw him overboard; suspended a cutlass over his naked wrists; Queequeg was the son of a King, and Queequeg budged not. Struck by his desperate dauntlessness, and his wild desire to visit Christendom, the captain at last relented, and told him he might make himself at home. But this fine young savage&mdash;this sea Prince of Wales, never saw the Captain's cabin. They put him down among the sailors, and made a whaleman of him. But like Czar Peter content to toil in the shipyards of foreign cities, Queequeg disdained no seeming ignominy, if thereby he might happily gain the power of enlightening his untutored countrymen. For at bottom&mdash;so he told me&mdash;he was actuated by a profound desire to learn among the Christians, the arts whereby to make his people still happier than they were; and more than that, still better than they were. But, alas! the practices of whalemen soon convinced him that even Christians could be both miserable and wicked; infinitely more so, than all his father's heathens. Arrived at last in old Sag Harbor; and seeing what the sailors did there; and then going on to Nantucket, and seeing how they spent their wages in that place also, poor Queequeg gave it up for lost. Thought he, it's a wicked world in all meridians; I'll die a pagan.</p>\r\n\r\n<p>And thus an old idolator at heart, he yet lived among these Christians, wore their clothes, and tried to talk their gibberish. Hence the queer ways about him, though now some time from home.</p>"
		a.Tags = []string{"Herman Melville", "Classics"}
		a.Save(revmgo.Session)

		a.Id = bson.ObjectIdHex("51e9ad9749a1b71843000002")
		a.Author_id = bson.ObjectIdHex("51e9aa4049a1b716bb000003")
		a.Published = true
		a.Posted = time.Date(2013, time.July, 18, 12, 32, 00, 0, time.UTC)
		a.Title = "Around the World in 80 Days"
		a.Body = "<p>\"The owners are myself,\" replied the captain.  \"The vessel belongs to me.\"</p>\r\n\r\n<p>\"I will freight it for you.\"</p>\r\n\r\n<p>\"No.\"</p>\r\n\r\n<p>\"I will buy it of you.\"</p>\r\n\r\n<p>\"No.\"</p>\r\n\r\n<p>Phileas Fogg did not betray the least disappointment; but the situation was a grave one.  It was not at New York as at Hong Kong, nor with the captain of the Henrietta as with the captain of the Tankadere.  Up to this time money had smoothed away every obstacle.  Now money failed.</p>\r\n\r\n<p>Still, some means must be found to cross the Atlantic on a boat, unless by balloon&mdash;which would have been venturesome, besides not being capable of being put in practice.  It seemed that Phileas Fogg had an idea, for he said to the captain, \"Well, will you carry me to Bordeaux?\"</p>\r\n\r\n<p>\"No, not if you paid me two hundred dollars.\"</p>"
		a.Tags = []string{"Jules Verne", "Classics", "Contemporary", "Action", "Adventure", "Suspense", "Fantasy"}
		a.Save(revmgo.Session)

		a.Id = bson.ObjectIdHex("51e9ae1749a1b71843000004")
		a.Author_id = bson.ObjectIdHex("51e9aa2d49a1b716bb000002")
		a.Published = true
		a.Posted = time.Date(2013, time.July, 17, 12, 32, 00, 0, time.UTC)
		a.Title = "A Princess of Mars"
		a.Body = "<p>Tal Hajus arose, and I, half fearing, half anticipating his intentions, hurried to the winding runway which led to the floors below.  No one was near to intercept me, and I reached the main floor of the chamber unobserved, taking my station in the shadow of the same column that Tars Tarkas had but just deserted.  As I reached the floor Tal Hajus was speaking.</p>\r\n\r\n<p>\"Princess of Helium, I might wring a mighty ransom from your people would I but return you to them unharmed, but a thousand times rather would I watch that beautiful face writhe in the agony of torture; it shall be long drawn out, that I promise you; ten days of pleasure were all too short to show the love I harbor for your race.  The terrors of your death shall haunt the slumbers of the red men through all the ages to come; they will shudder in the shadows of the night as their fathers tell them of the awful vengeance of the green men; of the power and might and hate and cruelty of Tal Hajus.  But before the torture you shall be mine for one short hour, and word of that too shall go forth to Tardos Mors, Jeddak of Helium, your grandfather, that he may grovel upon the ground in the agony of his sorrow.  Tomorrow the torture will commence; tonight thou art Tal Hajus'; come!\"</p>\r\n\r\n<p>He sprang down from the platform and grasped her roughly by the arm, but scarcely had he touched her than I leaped between them.  My short-sword, sharp and gleaming was in my right hand; I could have plunged it into his putrid heart before he realized that I was upon him; but as I raised my arm to strike I thought of Tars Tarkas, and, with all my rage, with all my hatred, I could not rob him of that sweet moment for which he had lived and hoped all these long, weary years, and so, instead, I swung my good right fist full upon the point of his jaw.  Without a sound he slipped to the floor as one dead.</p>\r\n\r\n<p>In the same deathly silence I grasped Dejah Thoris by the hand, and motioning Sola to follow we sped noiselessly from the chamber and to the floor above.  Unseen we reached a rear window and with the straps and leather of my trappings I lowered, first Sola and then Dejah Thoris to the ground below.  Dropping lightly after them I drew them rapidly around the court in the shadows of the buildings, and thus we returned over the same course I had so recently followed from the distant boundary of the city.</p>\r\n\r\n<p>We finally came upon my thoats in the courtyard where I had left them, and placing the trappings upon them we hastened through the building to the avenue beyond.  Mounting, Sola upon one beast, and Dejah Thoris behind me upon the other, we rode from the city of Thark through the hills to the south.</p>\r\n\r\n<p>Instead of circling back around the city to the northwest and toward the nearest waterway which lay so short a distance from us, we turned to the northeast and struck out upon the mossy waste across which, for two hundred dangerous and weary miles, lay another main artery leading to Helium.</p>"
		a.Tags = []string{"Edgar Rice Burroughs", "Adventure"}
		a.Save(revmgo.Session)

		a.Id = bson.ObjectIdHex("51e9ae4949a1b71843000005")
		a.Author_id = bson.ObjectIdHex("51e9aa2d49a1b716bb000002")
		a.Published = true
		a.Posted = time.Date(2013, time.July, 16, 12, 32, 00, 0, time.UTC)
		a.Title = "At the Earth's Core"
		a.Body = "<p>With no heavenly guide, it is little wonder that I became confused and lost in the labyrinthine maze of those mighty hills.  What, in reality, I did was to pass entirely through them and come out above the valley upon the farther side.  I know that I wandered for a long time, until tired and hungry I came upon a small cave in the face of the limestone formation which had taken the place of the granite farther back.</p>\r\n\r\n<p>The cave which took my fancy lay halfway up the precipitous side of a lofty cliff.  The way to it was such that I knew no extremely formidable beast could frequent it, nor was it large enough to make a comfortable habitat for any but the smaller mammals or reptiles.  Yet it was with the utmost caution that I crawled within its dark interior.</p>\r\n\r\n<p>Here I found a rather large chamber, lighted by a narrow cleft in the rock above which let the sunlight filter in in sufficient quantities partially to dispel the utter darkness which I had expected.  The cave was entirely empty, nor were there any signs of its having been recently occupied.  The opening was comparatively small, so that after considerable effort I was able to lug up a bowlder from the valley below which entirely blocked it.</p>\r\n\r\n<p>Then I returned again to the valley for an armful of grasses and on this trip was fortunate enough to knock over an orthopi, the diminutive horse of Pellucidar, a little animal about the size of a fox terrier, which abounds in all parts of the inner world.  Thus, with food and bedding I returned to my lair, where after a meal of raw meat, to which I had now become quite accustomed, I dragged the bowlder before the entrance and curled myself upon a bed of grasses&mdash;a naked, primeval, cave man, as savagely primitive as my prehistoric progenitors.</p>\r\n\r\n<p>I awoke rested but hungry, and pushing the bowlder aside crawled out upon the little rocky shelf which was my front porch.  Before me spread a small but beautiful valley, through the center of which a clear and sparkling river wound its way down to an inland sea, the blue waters of which were just visible between the two mountain ranges which embraced this little paradise.  The sides of the opposite hills were green with verdure, for a great forest clothed them to the foot of the red and yellow and copper green of the towering crags which formed their summit.  The valley itself was carpeted with a luxuriant grass, while here and there patches of wild flowers made great splashes of vivid color against the prevailing green.</p>"
		a.Tags = []string{"Edgar Rice Burroughs", "Adventure", "Action", "Fantasy", "Science Fiction"}
		a.Save(revmgo.Session)

		a.Id = bson.ObjectIdHex("51e9af2749a1b71843000006")
		a.Author_id = bson.ObjectIdHex("51e9aa2d49a1b716bb000002")
		a.Published = true
		a.Posted = time.Date(2013, time.July, 15, 12, 32, 00, 0, time.UTC)
		a.Title = "The War of the Worlds Book I"
		a.Body = "<p>\"Did you see a man in the pit?\" I said; but he made no answer to that.  We became silent, and stood watching for a time side by side, deriving, I fancy, a certain comfort in one another's company.  Then I shifted my position to a little knoll that gave me the advantage of a yard or more of elevation and when I looked for him presently he was walking towards Woking.</p>\r\n\r\n<p>The sunset faded to twilight before anything further happened.  The crowd far away on the left, towards Woking, seemed to grow, and I heard now a faint murmur from it.  The little knot of people towards Chobham dispersed.  There was scarcely an intimation of movement from the pit.</p>\r\n\r\n<p>It was this, as much as anything, that gave people courage, and I suppose the new arrivals from Woking also helped to restore confidence.  At any rate, as the dusk came on a slow, intermittent movement upon the sand pits began, a movement that seemed to gather force as the stillness of the evening about the cylinder remained unbroken.  Vertical black figures in twos and threes would advance, stop, watch, and advance again, spreading out as they did so in a thin irregular crescent that promised to enclose the pit in its attenuated horns.  I, too, on my side began to move towards the pit.</p>\r\n\r\n<p>Then I saw some cabmen and others had walked boldly into the sand pits, and heard the clatter of hoofs and the gride of wheels.  I saw a lad trundling off the barrow of apples.  And then, within thirty yards of the pit, advancing from the direction of Horsell, I noted a little black knot of men, the foremost of whom was waving a white flag.</p>\r\n\r\n<p>This was the Deputation.  There had been a hasty consultation, and since the Martians were evidently, in spite of their repulsive forms, intelligent creatures, it had been resolved to show them, by approaching them with signals, that we too were intelligent.</p>"
		a.Tags = []string{"H. G. Wells", "Science Fiction", "Classics"}
		a.Save(revmgo.Session)
	}
}
