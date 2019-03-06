package configo

import (
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestFirstOrNopFixture(t *testing.T) {
	gunit.Run(new(FirstOrNopFixture), t)
}

type FirstOrNopFixture struct {
	*gunit.Fixture

	source1 *FakeSource
	source2 *FakeSource
}

func (this *FirstOrNopFixture) Setup() {
	this.source1 = new(FakeSource)
	this.source2 = new(FakeSource)
}

func (this *FirstOrNopFixture) TestFirstOfTwoNonNilSourcesUsed() {
	this.So(FirstOrNop(this.source1, this.source2), should.Equal, this.source1)
}

func (this *FirstOrNopFixture) TestSecondChosenWhenFirstIsUntypedNil() {
	this.So(FirstOrNop(nil, this.source2), should.Equal, this.source2)
}
func (this *FirstOrNopFixture) TestSecondChosenWhenFirstIsTypedNil() {
	this.source1 = nil
	this.So(FirstOrNop(this.source1, this.source2), should.Equal, this.source2)
}
