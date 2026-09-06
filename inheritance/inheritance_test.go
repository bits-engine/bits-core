package inheritance_test

import (
	"testing"

	"github.com/bits-engine/bits-core/inheritance"
	"github.com/bits-engine/bits-ecs/component"
	"github.com/bits-engine/bits-ecs/world"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInheritance_AddRemove(t *testing.T) {
	w := world.New(&world.Conf{WorkersCount: 2})

	inheritanceTypes := inheritance.Register(w)

	entParent := w.New()
	entChild1 := w.New()
	entChild2 := w.New()
	inheritance.Set(w.CS(), inheritanceTypes, entParent, entChild1)
	inheritance.Set(w.CS(), inheritanceTypes, entParent, entChild2)

	children, exists := component.Get(w.CS(), entParent, inheritanceTypes.Children)
	require.True(t, exists, "Children component does not exists on parent")
	assert.Contains(t, children.List(), entChild1)
	assert.Contains(t, children.List(), entChild2)

	parent1, exists := component.Get(w.CS(), entChild1, inheritanceTypes.Parent)
	require.True(t, exists, "Parent component does not exists on child 1")
	assert.Equal(t, entParent, parent1.Entity())

	parent2, exists := component.Get(w.CS(), entChild2, inheritanceTypes.Parent)
	require.True(t, exists, "Parent component does not exists on child 2")
	assert.Equal(t, entParent, parent2.Entity())

	// Removed 1
	inheritance.Remove(w.CS(), inheritanceTypes, entParent, entChild1)

	children, exists = component.Get(w.CS(), entParent, inheritanceTypes.Children)
	require.True(t, exists, "Children component does not exists on parent")
	assert.NotContains(t, children.List(), entChild1)
	assert.Contains(t, children.List(), entChild2)

	parent1, exists = component.Get(w.CS(), entChild1, inheritanceTypes.Parent)
	assert.False(t, exists, "Parent component exists on child 1 after deletion")

	parent2, exists = component.Get(w.CS(), entChild2, inheritanceTypes.Parent)
	require.True(t, exists, "Parent component does not exists on child 2")
	assert.Equal(t, entParent, parent2.Entity())

	// Removed 2
	inheritance.Remove(w.CS(), inheritanceTypes, entParent, entChild2)

	children, exists = component.Get(w.CS(), entParent, inheritanceTypes.Children)
	assert.False(t, exists, "Children component exists on parent after deleting all children")

	parent1, exists = component.Get(w.CS(), entChild1, inheritanceTypes.Parent)
	assert.False(t, exists, "Parent component exists on child 1 after deletion")

	parent2, exists = component.Get(w.CS(), entChild2, inheritanceTypes.Parent)
	assert.False(t, exists, "Parent component exists on child 2 after deletion")
}

func TestInheritance_ParentChange(t *testing.T) {
	w := world.New(&world.Conf{WorkersCount: 2})

	inheritanceTypes := inheritance.Register(w)

	entParent1 := w.New()
	entParent2 := w.New()
	entChild := w.New()

	inheritance.Set(w.CS(), inheritanceTypes, entParent1, entChild)

	children, exists := component.Get(w.CS(), entParent1, inheritanceTypes.Children)
	require.True(t, exists, "Children component does not exists on parent 1")
	assert.Contains(t, children.List(), entChild)

	children, exists = component.Get(w.CS(), entParent2, inheritanceTypes.Children)
	assert.False(t, exists, "Children component exists on parent 2")

	parent, exists := component.Get(w.CS(), entChild, inheritanceTypes.Parent)
	require.True(t, exists, "Parent component does not exists on child")
	assert.Equal(t, entParent1, parent.Entity())

	// Change parent
	inheritance.Set(w.CS(), inheritanceTypes, entParent2, entChild)

	children, exists = component.Get(w.CS(), entParent2, inheritanceTypes.Children)
	require.True(t, exists, "Children component does not exists on parent 2")
	assert.Contains(t, children.List(), entChild)

	children, exists = component.Get(w.CS(), entParent1, inheritanceTypes.Children)
	assert.False(t, exists, "Children component exists on parent 1")

	parent, exists = component.Get(w.CS(), entChild, inheritanceTypes.Parent)
	require.True(t, exists, "Parent component does not exists on child")
	assert.Equal(t, entParent2, parent.Entity())
}
