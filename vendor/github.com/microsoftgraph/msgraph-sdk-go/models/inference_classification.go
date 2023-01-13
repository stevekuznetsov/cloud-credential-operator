package models

import (
    i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f "github.com/microsoft/kiota-abstractions-go"
    i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91 "github.com/microsoft/kiota-abstractions-go/serialization"
)

// InferenceClassification 
type InferenceClassification struct {
    Entity
    // A set of overrides for a user to always classify messages from specific senders in certain ways: focused, or other. Read-only. Nullable.
    overrides []InferenceClassificationOverrideable
}
// NewInferenceClassification instantiates a new inferenceClassification and sets the default values.
func NewInferenceClassification()(*InferenceClassification) {
    m := &InferenceClassification{
        Entity: *NewEntity(),
    }
    odataTypeValue := "#microsoft.graph.inferenceClassification";
    m.SetOdataType(&odataTypeValue);
    return m
}
// CreateInferenceClassificationFromDiscriminatorValue creates a new instance of the appropriate class based on discriminator value
func CreateInferenceClassificationFromDiscriminatorValue(parseNode i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable, error) {
    return NewInferenceClassification(), nil
}
// GetFieldDeserializers the deserialization information for the current model
func (m *InferenceClassification) GetFieldDeserializers()(map[string]func(i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.ParseNode)(error)) {
    res := m.Entity.GetFieldDeserializers()
    res["overrides"] = i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.SetCollectionOfObjectValues(CreateInferenceClassificationOverrideFromDiscriminatorValue , m.SetOverrides)
    return res
}
// GetOverrides gets the overrides property value. A set of overrides for a user to always classify messages from specific senders in certain ways: focused, or other. Read-only. Nullable.
func (m *InferenceClassification) GetOverrides()([]InferenceClassificationOverrideable) {
    return m.overrides
}
// Serialize serializes information the current object
func (m *InferenceClassification) Serialize(writer i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.SerializationWriter)(error) {
    err := m.Entity.Serialize(writer)
    if err != nil {
        return err
    }
    if m.GetOverrides() != nil {
        cast := i2ae4187f7daee263371cb1c977df639813ab50ffa529013b7437480d1ec0158f.CollectionCast[i878a80d2330e89d26896388a3f487eef27b0a0e6c010c493bf80be1452208f91.Parsable](m.GetOverrides())
        err = writer.WriteCollectionOfObjectValues("overrides", cast)
        if err != nil {
            return err
        }
    }
    return nil
}
// SetOverrides sets the overrides property value. A set of overrides for a user to always classify messages from specific senders in certain ways: focused, or other. Read-only. Nullable.
func (m *InferenceClassification) SetOverrides(value []InferenceClassificationOverrideable)() {
    m.overrides = value
}